package dbclient

type InsertBuilder struct {
	query string
	data  interface{}
}

func (c *DBClient) NewInsertBuilder() *InsertBuilder { return &InsertBuilder{} }

func (i *InsertBuilder) Query(query string) *InsertBuilder {
	i.query = query
	return i
}

func (i *InsertBuilder) Data(data interface{}) *InsertBuilder {
	i.data = data
	return i
}
