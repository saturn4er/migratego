package migratego

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type queryBuilderFunc func(QueryBuilder)

type migrationExecutor interface {
	Up(QueryBuilder)
	Down(QueryBuilder)
}
type Migration struct{}

func (m *Migration) Up(QueryBuilder)   {}
func (m *Migration) Down(QueryBuilder) {}

type MigrateApplication interface {
	AddMigration(int, string, migrationExecutor)
	SetSchemaVersionTable(string)
	Run([]string)
}

type migrateApplication struct {
	driver         string
	dsn            string
	dbVersionTable string
	migrations     []DBMigration
	db             *sql.DB
}

func (m *migrateApplication) AddMigration(number int, name string, me migrationExecutor) {
	for _, mi := range m.migrations {
		if mi.Number == number {
			fmt.Println("Error while adding migration " + name + ": migration with such number already exists")
			os.Exit(1)
		}
	}
	upScripts := m.getQueryBuilderScripts(me.Up)
	downScripts := m.getQueryBuilderScripts(me.Down)
	m.migrations = append(m.migrations, DBMigration{
		Name:       name,
		Number:     number,
		UpScript:   strings.Join(upScripts, ";"),
		DownScript: strings.Join(downScripts, ";"),
	})
	reflect.TypeOf(func(ab string) {}).String()
}
func (m *migrateApplication) SetSchemaVersionTable(name string) {
	m.dbVersionTable = name
}
func (m *migrateApplication) Run(args []string) {
	err := RunToolCli(m, args)
	if err != nil {
		fmt.Println(err)
	}
}
func (m *migrateApplication) getQueryBuilderScripts(p queryBuilderFunc) []string {
	qb := getDriverQueryBuilder(m.driver)
	p(qb)
	return qb.Sqls()
}
func (m *migrateApplication) getDriverClient() (DBClient, error) {
	return getDriverClient(m.driver, m.dsn, m.dbVersionTable)
}

func NewApp(driver, dsn string) MigrateApplication {
	shouldCheckDriver(driver)
	result := new(migrateApplication)
	result.dsn = dsn
	result.driver = driver
	result.dbVersionTable = "schema_version"
	return result
}
