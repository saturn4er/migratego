package migratego

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/saturn4er/migratego/databases/mysql"
	"github.com/saturn4er/migratego/types"
)

type QueryBuilder interface {
	DropTable(...string) types.DropTablesGenerator
	CreateTable(string, func(types.CreateTableGenerator)) types.CreateTableGenerator
	RawQuery(string)
	GetSqls() []string
}

type queryBuilderFunc func(QueryBuilder)

type MigrateApplication interface {
	AddMigration(int, string, queryBuilderFunc, queryBuilderFunc)
	SetSchemaVersionTable(string)
	Run([]string)
}

type migrateApplication struct {
	driver         string
	dsn            string
	dbVersionTable string
	migrations     []types.Migration
	db             *sql.DB
}

func (m *migrateApplication) AddMigration(number int, name string, up queryBuilderFunc, down queryBuilderFunc) {
	for _, mi := range m.migrations {
		if mi.Number == number {
			fmt.Println("Error while adding migration " + name + ": migration with such number already exists")
			os.Exit(1)
		}
	}
	upScripts := m.getQueryBuilderScripts(up)
	downScripts := m.getQueryBuilderScripts(down)
	m.migrations = append(m.migrations, types.Migration{
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
	return qb.GetSqls()
}
func getDriverQueryBuilder(driver string) QueryBuilder {
	switch driver {
	case "mysql":
		return new(mysql.MysqlQueryBuilder)
	}
	panic("Unknown driver:" + driver)
}
func getDriverClient(driver, dsn, transactionsTableName string) types.Client {
	switch driver {
	case "mysql":
		return mysql.NewClient(dsn, transactionsTableName)
	}
	panic("Unknown driver:" + driver)
}
func NewApp(driver, dsn string) MigrateApplication {
	shouldCheckDriver(driver)
	result := new(migrateApplication)
	result.dsn = dsn
	result.driver = driver
	result.dbVersionTable = "shema_version"
	return result
}