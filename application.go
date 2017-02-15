package migrates

import (
	"errors"

	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"
)

var commands []cli.Command

type MigrateApplication struct {
	dsn                string
	schemaVersionTable string
	migrations         []Migration
}

func (m *MigrateApplication) AddMigration(version int, name string, up func(*Scope), down func(*Scope)) error {
	for _, mi := range m.migrations {
		if mi.Version == version {
			return errors.New("")
		}
	}
	upScope := &Scope{}
	up(upScope)
	upScripts := make([]string, len(upScope.queries))
	for i, q := range upScope.queries {
		upScripts[i] = q.Sql()
	}
	downScope := &Scope{}
	down(downScope)
	downScripts := make([]string, len(downScope.queries))
	for i, q := range downScope.queries {
		downScripts[i] = q.Sql()
	}
	m.migrations = append(m.migrations, Migration{
		Name:       name,
		Version:    version,
		UpScript:   strings.Join(upScripts, ";"),
		DownScript: strings.Join(downScripts, ";"),
	})

	return nil
}
func (m *MigrateApplication) SetSchemaVersionTable(name string) {
	m.schemaVersionTable = name
}
func (m *MigrateApplication) Run(args []string) {
	app := cli.NewApp()
	app.Version = "1.0"
	app.Commands = commands
	if app.Metadata == nil {
		app.Metadata = make(map[string]interface{})
	}
	app.Metadata["migrations"] = m.migrations
	app.Run(args)
}

func NewApp(dsn string) *MigrateApplication {
	sqlx.Connect("mysql", dsn)
	result := new(MigrateApplication)
	return result
}
