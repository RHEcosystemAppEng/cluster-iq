package repositories

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/jmoiron/sqlx"
)

// TODO: Refactor to a Query Builder pattern.
// The current string concatenation approach is brittle and error-prone, especially with complex queries involving GROUP BY.
// A Query Builder would provide a more robust and maintainable way to construct SQL queries programmatically.
func listQueryHelper(ctx context.Context, db *sqlx.DB, dest interface{}, baseQuery string, countQuery string, opts models.ListOptions, whereClauses []string, namedArgs map[string]interface{}) (int, error) {
	if len(whereClauses) > 0 {
		whereStr := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += whereStr
		countQuery += whereStr
	}

	// Counting
	countStmt, countArgs, err := sqlx.Named(countQuery, namedArgs)
	if err != nil {
		return 0, fmt.Errorf("failed to bind named args for count query: %w", err)
	}
	countStmt = db.Rebind(countStmt)

	var total int
	if err := db.GetContext(ctx, &total, countStmt, countArgs...); err != nil {
		return 0, fmt.Errorf("failed to execute count query: %w", err)
	}

	if total == 0 {
		return 0, nil
	}

	// Pagination
	if opts.PageSize > 0 {
		// Only add ORDER BY if not already present in the base query
		if !strings.Contains(strings.ToUpper(baseQuery), "ORDER BY") {
			baseQuery += " ORDER BY id ASC"
		}
		baseQuery += " LIMIT :pagesize OFFSET :offset"
		namedArgs["pagesize"] = opts.PageSize
		namedArgs["offset"] = opts.Offset
	}

	// Main select
	queryStmt, queryArgs, err := sqlx.Named(baseQuery, namedArgs)
	if err != nil {
		return 0, fmt.Errorf("failed to bind named args for select query: %w", err)
	}
	queryStmt = db.Rebind(queryStmt)

	if err := db.SelectContext(ctx, dest, queryStmt, queryArgs...); err != nil {
		baseErr := fmt.Errorf("failed to execute select query: %w", err)
		if os.Getenv("CIQ_LOG_LEVEL") == "DEBUG" {
			return 0, fmt.Errorf("%w; query: %s, args: %v", baseErr, queryStmt, queryArgs)
		}
		return 0, baseErr
	}

	return total, nil
}
