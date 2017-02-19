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
