package migrates

type Querier interface {
	Sql() string
}
type Scope struct {
	queries []Querier
}
type Query string

func (s Query) Sql() string {
	return string(s)
}
func (s *Scope) DropTable(tableNames ...string) DropTablesGenerator {
	g := NewDropTablesGenerator(tableNames...)
	s.queries = append(s.queries, g)
	return g
}

func (s *Scope) CreateTable(tableName string, sc func(scope CreateTable)) {
	g := NewCreateTableGenerator(tableName, sc)
	s.queries = append(s.queries, g)
}

func (s *Scope) RawQuery(query string) {
	s.queries = append(s.queries, Query(query))

}

//func (s *Scope) AlterTable(tableName string, sc func(scope CreateTableQuery)) {
//	query := NewTable(tableName, sc)
//	s.queries = append(s.queries, query)
//}
