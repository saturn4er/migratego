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
func (m *MysqlQueryBuilder) GetSqls() []string {
	var result []string
	for _, g := range m.generators {
		result = append(result, g.Sql())
	}
	return result
}
