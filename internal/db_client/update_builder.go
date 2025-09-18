package dbclient

type UpdateBuilder struct {
	query string
	data  interface{}
}

func (c *DBClient) NewUpdateBuilder() *UpdateBuilder { return &UpdateBuilder{} }

func (i *UpdateBuilder) Query(query string) *UpdateBuilder {
	i.query = query
	return i
}

func (i *UpdateBuilder) Data(data interface{}) *UpdateBuilder {
	i.data = data
	return i
}
