package migratego

type MigrateApplication interface {
	AddMigration(int, string, migrationExecutor)
	SetSchemaVersionTable(string)
	Run()
}

type migrationExecutor interface {
	Up(QueryBuilder)
	Down(QueryBuilder)
}

type Migration struct{}

func (m *Migration) Up(QueryBuilder)   {}
func (m *Migration) Down(QueryBuilder) {}
