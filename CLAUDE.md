# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project Overview

**ClusterIQ** is an inventory and cost estimation platform for OpenShift clusters across multi-cloud environments (currently AWS only). It provides automated discovery, cost tracking, and lifecycle management.

**Architecture Components:**
1. **Scanner**: CronJob that discovers cloud resources using "Stocker" pattern
2. **API Server**: REST API (Gin framework) for inventory queries and cluster operations
3. **Agent**: gRPC service handling cluster power operations (instant, scheduled, recurring)

**Repository Structure:**
```
cmd/                    # Application entry points (api, scanner, agent)
internal/
  ├── actions/          # Action system (instant, scheduled, cron)
  ├── cloud_providers/  # AWS, GCP, Azure abstractions
  ├── cloud_executors/  # Action execution
  ├── stocker/          # Resource discovery
  ├── repositories/     # Data access layer
  ├── services/         # Business logic
  ├── api/handlers/     # HTTP handlers
  └── models/           # DTO, DB, domain models
db/sql/                 # Schema definitions (init.sql, cron.sql)
test/integration/       # Integration tests
```

## Essential Commands

**CRITICAL**: Always use Makefile. Never use direct `go build`, `go test`, or `go run`.
This includes compilation checks after code changes — use `make local-build` instead of `go build ./...`.

```bash
# Development
make check-dependencies  # Check before building
make all                # Stop, build, start dev environment
make start-dev          # Start services
make stop-dev           # Stop services

# Building
make build              # All container images
make local-build        # All local binaries

# Testing
make go-unit-tests
make go-integration-tests
make go-tests           # All tests
make lint               # Lint entire project
make lint-staged        # Lint staged files only

# Code Generation
make generate-converters  # Goverter (DB to DTO)
make swagger-doc         # OpenAPI docs
```

## Architecture Patterns

**Layered Architecture:**
- Handlers → Services → Repositories → Database
- Each layer has interfaces for testability
- Repository returns `repositories.ErrNotFound` for missing resources
- Repository returns `repositories.ErrNoClustersInAccount` when an account exists but has no clusters
- Services wrap all errors with `fmt.Errorf` context before returning to handlers
- Handlers map errors to HTTP status codes (404, 400, 500)

**Key Patterns:**
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic with dependency injection
- **Stocker Pattern**: Cloud resource discovery via `MakeStock()`
- **Action Channel**: Unbuffered channel for decoupled action dispatch
- **Event Tracker**: Audit logging with `EventService.StartTracking()`
- **Functional Options**: AWS connections use `WithEC2()`, `WithRoute53()`, etc.

## Testing Guidelines

**Service Layer Mocking:**
```go
// Mock repositories with function pointers
type mockClusterRepository struct {
    getClusterByIDFn func(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error)
}

func (m *mockClusterRepository) GetClusterByID(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error) {
    if m.getClusterByIDFn != nil {
        return m.getClusterByIDFn(ctx, clusterID)
    }
    return nil, nil
}
```

**gRPC Mocking:**
- Mock the underlying `pb.AgentServiceClient` interface, not the wrapper
- Return actual Go errors: `return nil, errTest`
- DON'T use response error codes for Go errors: `&pb.Response{Error: 1}` returns `nil` error

**Test Organization:**
- Group by method: `TestServiceName_MethodName(t *testing.T)`
- Use sub-tests: `t.Run("success", func(t *testing.T) { ... })`
- Naming: `testServiceName_MethodName_Scenario`
- All interface methods must be implemented in mocks (even if unused)

**Coverage Analysis:**
```bash
go test -coverprofile=coverage.out ./internal/services/...
go tool cover -func=coverage.out | grep -v "100.0%"  # Find uncovered
go tool cover -html=coverage.out -o coverage.html    # Visual
```

## Development Workflow

1. Make code changes
2. Run `make lint-staged` before committing
3. Run relevant tests: `make go-unit-tests`
4. For API changes: update Swagger with `make swagger-doc`
5. For DB changes: update `db/sql/init.sql` or add data migration in `doc/releases/`
6. For protobuf changes: `make local-build-agent`
7. For goverter changes: `make generate-converters`

**Commit Convention:**
- Use conventional commits format: `type(scope): brief description`
- Keep messages to single line
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`
- Examples:
  - `feat(api): implement PATCH endpoint for accounts`
  - `fix(aws): continue processing instances on conversion error`
  - `test(dto): add error handling tests for converters`

## Critical Implementation Notes

**Cloud Support:**
- AWS: Fully implemented ✅
- GCP/Azure: Interface stubs only ❌

**Data Consistency:**
- Each repository method = separate transaction
- No distributed transactions
- Eventual consistency for action status updates

**PATCH Endpoints:**
- Use pointer fields in DTOs to distinguish null from empty
- Dynamic SQL query building with positional parameters
- Verify resource exists before updating
- Refresh materialized views after successful update
- Return updated resource with HTTP 200

## Configuration

Environment variables:
- `CIQ_AGENT_URL` - Agent gRPC endpoint (default: "agent:50051")
- `CIQ_API_URL` - API endpoint used by Scanner and Agent (not by the API server itself)
- `CIQ_DB_URL` - PostgreSQL connection (default: "postgresql://pgsql:5432/clusteriq")
- `CIQ_CREDS_FILE` - Cloud provider credentials file
- `CIQ_LOG_LEVEL` - Log verbosity (default: "INFO")

## AWS Integration Details

**Cluster Discovery:**
- OpenShift tags: `kubernetes.io/cluster/<cluster-name>-<infraID>`
- Console URLs via Route53 DNS: `console-openshift-console.apps.*`

**Instance Timestamps:**
- Primary: Root EBS volume `AttachTime` (via `instance.RootDeviceName`)
- Fallback: Instance `LaunchTime`
- NEVER hardcode device names (varies by instance type)

**Cost Explorer:**
- 14-day rolling window, DAILY granularity
- Metric: UnblendedCost
- One Expense object per instance per day

**EC2 Operations:**
- Start: Only targets `stopped` instances
- Stop: Only targets `running` instances
- Prevents idempotency errors

## Code Generation

**Goverter**: Type-safe DB ↔ DTO converters
- Definitions: `internal/models/convert/`
- Generated: `internal/models/convert/generated.go`

**Swagger**: OpenAPI from Go comments
- Annotations in handlers and `cmd/api/server.go`

**Protobuf**: gRPC code from `cmd/agent/proto/agent.proto`
- Auto-generated during `make local-build-agent`
