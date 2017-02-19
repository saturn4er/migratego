package types
type Order string
type Querier interface {
	Sql() string
}
type DropTablesGenerator interface {
	Table(name string) DropTablesGenerator
	IfExists() DropTablesGenerator
	Sql() string
}
type CreateTableGenerator interface {
	Column(name string, Type string) ColumnGenerator
	Index(name string, unique bool) IndexGenerator
	NewIndexColumn(ColumnGenerator, ...interface {}) IndexColumnGenerator
	Sql() string
}
type ColumnGenerator interface {
	GetName() string
	NotNull() ColumnGenerator
	Primary(comment ...string) ColumnGenerator
	Index(name string, unique bool, order string, length int) IndexGenerator
	Sql() string
}
type IndexGenerator interface {
	Columns(...IndexColumnGenerator) IndexGenerator
	Sql() string
}
type IndexColumnGenerator interface {
	Sql() string
}

type Client interface{
	PrepareTransactionsTable() error
}