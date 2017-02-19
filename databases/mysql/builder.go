package mysql

import "github.com/saturn4er/migratego/types"

type MysqlQueryBuilder struct {
	generators []types.Querier
}

func (m *MysqlQueryBuilder) DropTable(names ...string) types.DropTablesGenerator {
	c := NewDropTablesGenerator(names...)
	m.generators = append(m.generators, c)
	return c
}
func (m *MysqlQueryBuilder) CreateTable(name string, g func(generator types.CreateTableGenerator)) types.CreateTableGenerator {
	c := NewCreateTableGenerator(name, g)
	m.generators = append(m.generators, c)
	return c
}
func (m *MysqlQueryBuilder) RawQuery(q string) {
	c := rawQuery(q)
	m.generators = append(m.generators, &c)
}
// NewIndexColumn creates new IndexColumnGenerator
// Usage NewIndexColumn(column, order[optional], length[optional])
// orderType default value is ASC
// length default value is int
func (c *MysqlQueryBuilder) NewIndexColumn(column types.ColumnGenerator, params ...interface{}) types.IndexColumnGenerator {
	var length int
	var order = "ASC"
	var ok bool
	if len(params) > 0 {
		if order, ok = params[0].(string); !ok {
			panic("first param should be of type `string`")
		}
	}
	if len(params) > 1 {
		if length, ok = params[0].(int); !ok {
			panic("first param should be of type `int`")
		}
	}
	return &IndexColumnGenerator{
		Column: column,
		Order:  order,
		Length: length,
	}
}
func (m *MysqlQueryBuilder) Sqls() []string {
	var result []string
	for _, g := range m.generators {
		result = append(result, g.Sql())
	}
	return result
}
