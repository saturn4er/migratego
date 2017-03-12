package migratego

type QueryBuilder interface {
	DropTables(...string) DropTablesGenerator
	CreateTable(string, func(CreateTableGenerator)) CreateTableGenerator
	Table(string, func(generator TableScope))
	NewIndexColumn(column string, params ...interface{}) IndexColumnGenerator
	RawQuery(string)
	Sqls() []string
}

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
	Column(name string, Type string) CreateTableColumnGenerator
	Index(name string, unique bool) IndexGenerator
	Engine(engine string) CreateTableGenerator
	Charset(charset string) CreateTableGenerator
	Comment(comment string) CreateTableGenerator
	Sql() string
}
type CreateTableColumnGenerator interface {
	GetName() string
	Primary(comment ...string) CreateTableColumnGenerator
	NotNull() CreateTableColumnGenerator
	Unsigned() CreateTableColumnGenerator
	Binary() CreateTableColumnGenerator
	ZeroFill() CreateTableColumnGenerator
	Generated() CreateTableColumnGenerator
	DefaultValue(v string) CreateTableColumnGenerator
	Comment(c string) CreateTableColumnGenerator
	AutoIncrement(primaryComment ...string) CreateTableColumnGenerator
	Index(name string, unique bool, params ...interface{}) IndexGenerator
	Sql() string
}
type UpdateTableGenerator interface {
	Rename(name string)
	Delete(name string)
	AddColumn(name string, Type string) UpdateTableAddColumnGenerator
	AddIndex(name string, unique bool) IndexGenerator
	RemoveIndex(name string)
	Sql() string
}
type UpdateTableAddColumnGenerator interface {
	GetName() string
	NotNull() UpdateTableAddColumnGenerator
	Unsigned() UpdateTableAddColumnGenerator
	Binary() UpdateTableAddColumnGenerator
	ZeroFill() UpdateTableAddColumnGenerator
	Generated() UpdateTableAddColumnGenerator
	DefaultValue(v string) UpdateTableAddColumnGenerator
	Comment(c string) UpdateTableAddColumnGenerator
	After(column string) UpdateTableAddColumnGenerator
	Sql() string
}
type IndexGenerator interface {
	Unique() IndexGenerator
	Columns(...IndexColumnGenerator) IndexGenerator
	Sql() string
}
type IndexColumnGenerator interface {
	Sql() string
}

type TableScope interface {
	AddColumn(name, cType string) UpdateTableAddColumnGenerator
	RemoveColumn(name string) TableScope
	AddIndex(string, bool) IndexGenerator
	RemoveIndex(name string) TableScope
	Rename(name string) TableScope
	Delete()
}
type DBClient interface {
	// PrepareTransactionsTable checks if table with migrations exists and creates it, if it doesn't
	PrepareTransactionsTable() error
	// Backup dumps database to some file in folder and returns path to it
	Backup(path string) (string, error)
	// InsertMigration adds migration to migrations table
	InsertMigration(migration *Migration) error
	// RemoveMigration removes migration from migrations table
	RemoveMigration(migration *Migration) error
	// ApplyMigration executes UpScript if down is false. Execute DownScript of down is true
	ApplyMigration(migration *Migration, down bool) error
	// GetAppliedMigrations returns list of migrations in migrations table
	GetAppliedMigrations() ([]Migration, error)
}
