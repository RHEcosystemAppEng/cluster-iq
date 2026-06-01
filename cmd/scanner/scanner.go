package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/scanner"
	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/credentials"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/stocker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	apiInventoryEndpoint = "/inventory"
	apiAccountEndpoint   = "/accounts"
	apiClusterEndpoint   = "/clusters"
	apiInstanceEndpoint  = "/instances"
	apiExpenseEndpoint   = "/expenses"
	// apiRequestTimeout defines the timeout for HTTP POST requests to the API
	apiRequestTimeout = 60 * time.Second
	// apiHealthcheckTimeout defines the timeout for each healthcheck attempt
	apiHealthcheckTimeout = 5 * time.Second
	// apiHealthcheckRetryInterval defines how long to wait between healthcheck retries
	apiHealthcheckRetryInterval = 3 * time.Second
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string

	// logger variable across the entire scanner code
	logger *zap.Logger

	// MD5 Checksum of credsFile
	credsFileHash []byte

	// HTTP Client for connecting the scanner to the API
	client http.Client

	// API URL as global var for postData function
	APIURL string
)

// Scanner models the cloud agnostic Scanner for looking up OCP deployments
type Scanner struct {
	pb.UnimplementedScannerServiceServer
	allAccounts map[string]credentials.AccountConfig
	cfg         *config.ScannerConfig
	logger      *zap.Logger
	grpcServer  *grpc.Server
	mu          sync.Mutex
	scanning    bool
}

// NewScanner creates and returns a new Scanner instance
func NewScanner(cfg *config.ScannerConfig, logger *zap.Logger) *Scanner {
	hash := md5.Sum([]byte(cfg.CredentialsFile))
	credsFileHash = hash[:]

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = http.Client{Transport: tr}
	APIURL = cfg.APIURL

	return &Scanner{
		allAccounts: make(map[string]credentials.AccountConfig),
		cfg:         cfg,
		logger:      logger,
	}
}

func init() {
	logger = ciqLogger.NewLogger()
}

// loadAccounts reads the credentials file and caches all account configs.
func (s *Scanner) loadAccounts() error {
	accountConfigs, err := credentials.ReadCloudAccounts(s.cfg.CredentialsFile)
	if err != nil {
		return fmt.Errorf("failed to read cloud accounts: %w", err)
	}

	for _, ac := range accountConfigs {
		s.allAccounts[ac.ID] = ac
	}

	s.logger.Info("Loaded cloud accounts from credentials file",
		zap.Int("count", len(s.allAccounts)))

	return nil
}

// waitForAPI blocks until the API server responds to /healthcheck.
func (s *Scanner) waitForAPI() {
	url := fmt.Sprintf("%s/healthcheck", APIURL)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), apiHealthcheckTimeout)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err == nil {
			resp, err := client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
			if err == nil && resp.StatusCode == http.StatusOK {
				cancel()
				s.logger.Info("API server is ready")
				return
			}
		}
		cancel()
		s.logger.Info("Waiting for API server...", zap.String("url", url))
		time.Sleep(apiHealthcheckRetryInterval)
	}
}

// seedAccounts posts account records from the credentials file to the API
// so they are available in the console before the first scan runs.
func (s *Scanner) seedAccounts() error {
	var accounts []dto.AccountDTORequest
	for _, ac := range s.allAccounts {
		accounts = append(accounts, dto.AccountDTORequest{
			AccountID:   ac.ID,
			AccountName: ac.Name,
			Provider:    ac.Provider,
		})
	}

	b, err := json.Marshal(accounts)
	if err != nil {
		return fmt.Errorf("failed to marshal seed accounts: %w", err)
	}

	if err := postData(apiAccountEndpoint, b); err != nil {
		return fmt.Errorf("failed to seed accounts: %w", err)
	}

	if err := refreshInventory(s.logger); err != nil {
		return fmt.Errorf("failed to refresh materialized views after seed: %w", err)
	}

	s.logger.Info("Seeded accounts into database", zap.Int("count", len(accounts)))
	return nil
}

// buildInventory creates an Inventory from the given account configs.
func (s *Scanner) buildInventory(configs []credentials.AccountConfig) (*inventory.Inventory, error) {
	inv := inventory.NewInventory()
	for _, ac := range configs {
		newAccount, err := inventory.NewAccount(ac.ID, ac.Name, ac.Provider, ac.User, ac.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to create account %s: %w", ac.ID, err)
		}
		if ac.BillingEnabled {
			newAccount.EnableBilling()
		}
		if err := inv.AddAccount(newAccount); err != nil {
			return nil, fmt.Errorf("failed to add account %s: %w", ac.ID, err)
		}
	}
	return inv, nil
}

// nolint:cyclop
func (s *Scanner) createStockers(inv *inventory.Inventory) ([]stocker.Stocker, []stocker.Stocker, []error) {
	var stockers []stocker.Stocker
	var billingStockers []stocker.Stocker
	var errs []error

	for _, account := range inv.Accounts {
		switch account.Provider {
		case inventory.AWSProvider:
			s.logger.Info("Processing AWS account", zap.String("account_id", account.AccountID), zap.String("account_name", account.AccountName))
			awsStocker, err := stocker.NewAWSStocker(account, s.cfg.SkipNoOpenShiftInstances, s.logger)
			if err != nil {
				s.logger.Error("Failed to create AWS stocker; skipping this account",
					zap.String("account", account.AccountName), zap.Error(err))
				errs = append(errs, fmt.Errorf("account %s: %w", account.AccountName, err))
				continue
			}
			stockers = append(stockers, awsStocker)

			if account.IsBillingEnabled() {
				s.logger.Info("Enabled AWS Billing Stocker", zap.String("account_id", account.AccountID), zap.String("account_name", account.AccountName))
				instancesToScan, err := getInstancesForBillingUpdate(s.cfg.APIURL, account.AccountID, s.logger)
				if err != nil {
					s.logger.Error("Failed to retrieve instances for billing",
						zap.String("account_name", account.AccountName), zap.Error(err))
				} else {
					if bs := stocker.NewAWSBillingStocker(account, s.logger, instancesToScan); bs != nil {
						billingStockers = append(billingStockers, bs)
					}
				}
			}
		case inventory.GCPProvider:
			s.logger.Warn("Failed to scan GCP account",
				zap.String("account_id", account.AccountID),
				zap.String("account_name", account.AccountName),
				zap.String("reason", "not implemented"))
		case inventory.AzureProvider:
			s.logger.Warn("Failed to scan Azure account",
				zap.String("account_id", account.AccountID),
				zap.String("account_name", account.AccountName),
				zap.String("reason", "not implemented"))
		case inventory.UnknownProvider:
			s.logger.Warn("Unknown cloud provider, skipping account",
				zap.String("account_id", account.AccountID),
				zap.String("account_name", account.AccountName),
				zap.String("provider", string(account.Provider)))
		default:
			s.logger.Warn("Unsupported cloud provider, skipping account",
				zap.String("account_id", account.AccountID),
				zap.String("account_name", account.AccountName),
				zap.String("provider", string(account.Provider)))
		}
	}

	s.logger.Info("Account registration complete",
		zap.Int("registeredAccounts", len(inv.Accounts)),
		zap.Int("registeredStockers", len(stockers)),
		zap.Int("skippedAccounts", len(inv.Accounts)-len(stockers)))

	return stockers, billingStockers, errs
}

func runStockers(stockers []stocker.Stocker, billingStockers []stocker.Stocker, l *zap.Logger) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(stockers)+len(billingStockers))

	l.Warn("Running Infrastructure Stockers!", zap.Int("stockers_count", len(stockers)))
	for _, st := range stockers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := st.MakeStock(); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()

	l.Warn("Running Billing Stockers!", zap.Int("stockers_count", len(billingStockers)))
	for _, st := range billingStockers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := st.MakeStock(); err != nil {
				errChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errorList []error
	for err := range errChan {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		for _, err := range errorList {
			l.Error("Stocker Error", zap.Error(err))
		}
		return fmt.Errorf("error when running Scanner stockers. Failed Stockers: (%d)", len(errorList))
	}

	l.Info("Stockers executed correctly")
	return nil
}

// selectAccounts returns account configs matching the given IDs, or all accounts if selectAll is true or no IDs are provided.
func (s *Scanner) selectAccounts(accountIDs []string, selectAll bool) []credentials.AccountConfig {
	var configs []credentials.AccountConfig
	if selectAll || len(accountIDs) == 0 {
		for _, ac := range s.allAccounts {
			configs = append(configs, ac)
		}
		return configs
	}

	for _, id := range accountIDs {
		ac, ok := s.allAccounts[id]
		if !ok {
			s.logger.Warn("Account not found in credentials, skipping", zap.String("account_id", id))
			continue
		}
		configs = append(configs, ac)
	}
	return configs
}

// ExecuteScan runs the full scan pipeline for the given accounts.
func (s *Scanner) ExecuteScan(accountIDs []string, selectAll bool) (int, error) {
	s.mu.Lock()
	if s.scanning {
		s.mu.Unlock()
		return 0, fmt.Errorf("a scan is already in progress")
	}
	s.scanning = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.scanning = false
		s.mu.Unlock()
	}()

	configs := s.selectAccounts(accountIDs, selectAll)
	if len(configs) == 0 {
		return 0, fmt.Errorf("no valid accounts found for scanning")
	}

	inv, err := s.buildInventory(configs)
	if err != nil {
		return 0, fmt.Errorf("failed to build inventory: %w", err)
	}

	stockers, billingStockers, stockerErrors := s.createStockers(inv)
	if len(stockers) == 0 {
		return 0, fmt.Errorf("no valid stockers created: %w", errors.Join(stockerErrors...))
	}

	if err := runStockers(stockers, billingStockers, s.logger); err != nil {
		return 0, fmt.Errorf("failed to run stockers: %w", err)
	}

	inv.PrintInventory()
	if err := postScannerInventory(inv, s.logger); err != nil {
		return 0, fmt.Errorf("failed to post inventory: %w", err)
	}

	return len(configs), nil
}

// Scan handles a gRPC scan request.
func (s *Scanner) Scan(_ context.Context, req *pb.ScanRequest) (*pb.ScanResponse, error) {
	s.logger.Info("Received scan request",
		zap.Int64("run_id", req.RunId),
		zap.Strings("account_ids", req.AccountIds),
		zap.Bool("select_all", req.SelectAll))

	scanned, err := s.ExecuteScan(req.AccountIds, req.SelectAll)
	if err != nil {
		s.logger.Error("Scan failed", zap.Error(err))
		return &pb.ScanResponse{
			Error:           1,
			Message:         err.Error(),
			AccountsScanned: 0,
		}, nil
	}

	s.logger.Info("Scan completed successfully", zap.Int("accounts_scanned", scanned))
	return &pb.ScanResponse{
		Error:           0,
		Message:         "Scan completed successfully",
		AccountsScanned: int32(scanned),
	}, nil
}

// Health reports the scanner's readiness.
func (s *Scanner) Health(_ context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Ready: true}, nil
}

// startGRPCServer initializes and starts the gRPC server.
func (s *Scanner) startGRPCServer() error {
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", s.cfg.ListenURL)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.ListenURL, err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterScannerServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	s.logger.Info("Scanner gRPC server listening", zap.String("address", s.cfg.ListenURL))

	return s.grpcServer.Serve(lis)
}

func postNewAccount(account inventory.Account, l *zap.Logger) error {
	l.Debug("Posting new Account", zap.String("account_id", account.AccountID), zap.String("account_name", account.AccountName))

	var accounts []dto.AccountDTORequest
	accounts = append(accounts, *dto.ToAccountDTORequest(account))
	b, err := json.Marshal(accounts)
	if err != nil {
		l.Error("Failed to marshal account", zap.String("account_id", account.AccountID), zap.Error(err))
		return err
	}

	if err := postData(apiAccountEndpoint, b); err != nil {
		return err
	}

	clusters, instances, expenses := flatternAccount(account)

	if len(clusters) > 0 {
		if err := postClusters(clusters); err != nil {
			return err
		}
	}

	if len(instances) > 0 {
		if err := postInstances(instances); err != nil {
			return err
		}
	}

	if len(expenses) > 0 {
		l.Info("Posting expenses", zap.Int("expenses_count", len(expenses)))
		if err := postExpenses(expenses); err != nil {
			return err
		}
	}
	return nil
}

func flatternAccount(account inventory.Account) ([]inventory.Cluster, []inventory.Instance, []inventory.Expense) {
	var clusters []inventory.Cluster
	var instances []inventory.Instance
	var expenses []inventory.Expense
	for _, cluster := range account.Clusters {
		for _, instance := range cluster.Instances {
			expenses = append(expenses, instance.Expenses...)
			instances = append(instances, instance)
		}
		clusters = append(clusters, *cluster)
	}
	return clusters, instances, expenses
}

func postClusters(clusters []inventory.Cluster) error {
	b, err := json.Marshal(dto.ToClusterDTORequestList(clusters))
	if err != nil {
		return err
	}
	return postData(apiClusterEndpoint, b)
}

func postInstances(instances []inventory.Instance) error {
	b, err := json.Marshal(dto.ToInstanceDTORequestList(instances))
	if err != nil {
		return err
	}
	return postData(apiInstanceEndpoint, b)
}

func postExpenses(expenses []inventory.Expense) error {
	b, err := json.Marshal(dto.ToExpenseDTORequestList(expenses))
	if err != nil {
		return err
	}
	return postData(apiExpenseEndpoint, b)
}

func postScannerInventory(inv *inventory.Inventory, l *zap.Logger) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(inv.Accounts))

	for _, account := range inv.Accounts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := postNewAccount(*account, l); err != nil {
				errChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errorList []error
	for err := range errChan {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		for _, err := range errorList {
			l.Error("Post Account Error", zap.Error(err))
		}
		return fmt.Errorf("error when posting Scanner inventory")
	}

	l.Info("Inventory posted correctly")

	if err := refreshInventory(l); err != nil {
		return err
	}

	return nil
}

func refreshInventory(l *zap.Logger) error {
	if err := postData(apiInventoryEndpoint, []byte{}); err != nil {
		return err
	}
	l.Info("Inventory refreshed")
	return nil
}

func postData(path string, b []byte) error {
	url := fmt.Sprintf("%s%s", APIURL, path)
	ctx, cancel := context.WithTimeout(context.Background(), apiRequestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("API returned HTTP %d for %s", response.StatusCode, path)
	}
	return nil
}

func getInstancesForBillingUpdate(apiURL string, accountID string, l *zap.Logger) ([]string, error) {
	l.Debug("Fetching instances for update billing from backend")

	requestURL := apiURL + apiAccountEndpoint + "/" + accountID + "/expense_update"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, requestURL, nil)
	if err != nil {
		l.Error("Failed preparing last expenses list request", zap.Error(err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Error("Failed to get last expenses from API", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Error("Failed to get last expenses from API", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to get last expenses, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Error("Failed to read response body", zap.Error(err))
		return nil, err
	}

	var response responsetypes.ListResponse[string]
	err = json.Unmarshal(body, &response)
	if err != nil {
		l.Error("Failed to unmarshal instance IDs JSON", zap.Error(err))
		return nil, err
	}

	l.Debug("Successfully fetched instance IDs from backend", zap.Int("instances_num", response.Count))
	return response.Items, nil
}

func main() {
	defer func() { _ = logger.Sync() }()

	cfg, err := config.LoadScannerConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	scan := NewScanner(cfg, logger)

	scan.logger.Info("==================== Starting ClusterIQ Scanner ====================",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("credentials_file_path", cfg.CredentialsFile),
		zap.ByteString("credentials_file_hash", credsFileHash),
		zap.String("listen_url", cfg.ListenURL),
	)

	if err := scan.loadAccounts(); err != nil {
		logger.Fatal("Failed to load cloud accounts", zap.Error(err))
	}

	scan.waitForAPI()
	if err := scan.seedAccounts(); err != nil {
		logger.Warn("Failed to seed accounts", zap.Error(err))
	}

	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-quit
		logger.Warn("Received signal, shutting down...", zap.String("signal", s.String()))
		if scan.grpcServer != nil {
			scan.grpcServer.GracefulStop()
		}
	}()

	if err := scan.startGRPCServer(); err != nil {
		logger.Fatal("gRPC server failed", zap.Error(err))
	}

	logger.Info("Scanner stopped")
}
