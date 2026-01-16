package dbclient

import (
	"fmt"
	"strings"
)

type DeleteBuilder struct {
	table string
	where []string
	args  []interface{}
}

func (d *DBClient) NewDeleteBuilder() *DeleteBuilder { return &DeleteBuilder{} }

func (d *DeleteBuilder) From(table string) *DeleteBuilder {
	d.table = table
	return d
}

func (d *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	d.where = append(d.where, condition)
	d.args = append(d.args, args...)
	return d
}

func (d *DeleteBuilder) Build() (string, []interface{}, error) {
	if d.table == "" {
		return "", nil, fmt.Errorf("no table defined for DeleteBuilder")
	}

	query := fmt.Sprintf("DELETE FROM %s", d.table)

	if len(d.where) > 0 {
		query += " WHERE " + strings.Join(d.where, " AND ")
	}

	return query, d.args, nil
}
