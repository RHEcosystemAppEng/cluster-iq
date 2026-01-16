package dbclient

import (
	"fmt"
	"strings"
)

// SelectBuilder builds SELECT queries with pagination
type SelectBuilder struct {
	columns []string
	table   string
	where   []string
	args    []interface{}
	limit   int
	offset  int
	orderBy string
}

// NewSelectBuilder creates a new SelectBuilder
func (d *DBClient) NewSelectBuilder(columns ...string) *SelectBuilder {
	return &SelectBuilder{
		columns: columns,
		limit:   -1,
		offset:  -1,
	}
}

func (b *SelectBuilder) From(table string) *SelectBuilder {
	b.table = table
	return b
}

func (b *SelectBuilder) Where(condition string, args ...interface{}) *SelectBuilder {
	b.where = append(b.where, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *SelectBuilder) OrderBy(clause string) *SelectBuilder {
	b.orderBy = clause
	return b
}

func (b *SelectBuilder) Paginate(limit, offset int) *SelectBuilder {
	b.limit = limit
	b.offset = offset
	return b
}

func (b *SelectBuilder) Build() (string, []interface{}, error) {
	cols := "*"
	if len(b.columns) > 0 {
		cols = strings.Join(b.columns, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", cols, b.table)

	if len(b.where) > 0 {
		query += " WHERE " + strings.Join(b.where, " AND ")
	}

	isPaginating := b.limit >= 0
	if isPaginating && b.orderBy == "" {
		return "", nil, fmt.Errorf("enabling pagination requires an ORDER BY column to ensure consistency")
	}

	if b.orderBy != "" {
		query += " ORDER BY " + b.orderBy + " DESC"
	}

	if isPaginating {
		query += fmt.Sprintf(" LIMIT %d", b.limit)
		if b.offset >= 0 {
			query += fmt.Sprintf(" OFFSET %d", b.offset)
		}
	}

	return query, b.args, nil
}
