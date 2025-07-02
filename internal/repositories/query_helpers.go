package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func listQueryHelper(ctx context.Context, db *sqlx.DB, dest interface{}, baseQuery string, countQuery string, opts ListOptions, whereClauses []string, namedArgs map[string]interface{}) (int, error) {
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
	baseQuery += " ORDER BY id LIMIT :pagesize OFFSET :offset"
	namedArgs["pagesize"] = opts.PageSize
	namedArgs["offset"] = opts.Offset

	// Main select
	queryStmt, queryArgs, err := sqlx.Named(baseQuery, namedArgs)
	if err != nil {
		return 0, fmt.Errorf("failed to bind named args for select query: %w", err)
	}
	queryStmt = db.Rebind(queryStmt)

	if err := db.SelectContext(ctx, dest, queryStmt, queryArgs...); err != nil {
		return 0, fmt.Errorf("failed to execute select query: %w", err)
	}

	return total, nil
}
