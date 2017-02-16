package migrates

import (
	"reflect"
	"strings"

	"database/sql"
	"fmt"
	"os"
)

type MigrateApplication struct {
	dsn            string
	dbVersionTable string
	migrations     []Migration
	db             *sql.DB
}

func (m *MigrateApplication) AddMigration(number int, name string, up func(*Scope), down func(*Scope)) {
	for _, mi := range m.migrations {
		if mi.Number == number {
			fmt.Println("Error while adding migration " + name + ": migration with such number already exists")
			os.Exit(1)
		}
	}
	upScripts := getScopeScripts(up)
	downScripts := getScopeScripts(down)
	m.migrations = append(m.migrations, Migration{
		Name:       name,
		Number:     number,
		UpScript:   strings.Join(upScripts, ";"),
		DownScript: strings.Join(downScripts, ";"),
	})
	reflect.TypeOf(func(ab string) {}).String()
}
func (m *MigrateApplication) SetSchemaVersionTable(name string) {
	m.dbVersionTable = name
}
func (m *MigrateApplication) Run(args []string) {
	err := RunToolCli(m, args)
	if err != nil {
		fmt.Println(err)
	}
}

func getScopeScripts(p func(*Scope)) []string {
	scope := new(Scope)
	p(scope)
	scripts := make([]string, len(scope.queries))
	for i, q := range scope.queries {
		scripts[i] = q.Sql()
	}
	return scripts
}

func NewApp(dsn string) *MigrateApplication {
	result := new(MigrateApplication)
	result.dsn = dsn
	result.dbVersionTable = "shema_version"
	return result
}
