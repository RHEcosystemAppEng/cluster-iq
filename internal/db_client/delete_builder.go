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

	if len(d.where) == 0 {
		return "", nil, fmt.Errorf("DELETE requires at least one WHERE condition for safety")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", d.table, strings.Join(d.where, " AND "))

	return query, d.args, nil
}
