package clients

import (
	"context"

	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/scanner"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ScannerGRPCClient manages the gRPC client connection to the Scanner service.
type ScannerGRPCClient struct {
	Client pb.ScannerServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewScannerGRPCClient initializes and returns a new ScannerGRPCClient.
func NewScannerGRPCClient(scannerURL string, logger *zap.Logger) (*ScannerGRPCClient, error) {
	conn, err := grpc.NewClient(scannerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ScannerGRPCClient{
		Client: pb.NewScannerServiceClient(conn),
		conn:   conn,
		logger: logger,
	}, nil
}

// Close closes the underlying gRPC connection.
func (s *ScannerGRPCClient) Close() error {
	return s.conn.Close()
}

// Scan sends a scan request to the Scanner service.
func (s *ScannerGRPCClient) Scan(ctx context.Context, runID int64, accountIDs []string, selectAll bool) (*pb.ScanResponse, error) {
	req := &pb.ScanRequest{
		RunId:      runID,
		AccountIds: accountIDs,
		SelectAll:  selectAll,
	}

	s.logger.Info("Sending scan request to Scanner",
		zap.Int64("run_id", runID),
		zap.Strings("account_ids", accountIDs),
		zap.Bool("select_all", selectAll),
	)

	resp, err := s.Client.Scan(ctx, req)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Scan response received",
		zap.Int32("error_code", resp.Error),
		zap.String("message", resp.Message),
		zap.Int32("accounts_scanned", resp.AccountsScanned),
	)

	return resp, nil
}
