package migrates

import (
	"errors"
	"reflect"
	"strings"

	"database/sql"
	"fmt"
)

type MigrateApplication struct {
	dsn            string
	dbVersionTable string
	migrations     []Migration
	db             *sql.DB
}

func (m *MigrateApplication) AddMigration(version int, name string, up func(*Scope), down func(*Scope)) error {
	for _, mi := range m.migrations {
		if mi.Version == version {
			return errors.New("Error while adding migration: with such version already exists")
		}
	}
	upScripts := getScopeScripts(up)
	downScripts := getScopeScripts(down)
	m.migrations = append(m.migrations, Migration{
		Name:       name,
		Version:    version,
		UpScript:   strings.Join(upScripts, ";"),
		DownScript: strings.Join(downScripts, ";"),
	})
	reflect.TypeOf(func(ab string) {}).String()

	return nil
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
