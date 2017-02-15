package migrates

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"fmt"
	"strings"
)

type MigrateApplication struct {
	Dsn        string
	Migrations []Migration
}

func (m *MigrateApplication) AddMigration(version int, name string, up func(*Scope), down func(*Scope)) error {
	for _, mi := range m.Migrations {
		if mi.Version == version {
			return errors.New("")
		}
	}
	upScope := &Scope{}
	up(upScope)
	upScripts := make([]string, len(upScope.queries))
	for i, q := range upScope.queries {
		upScripts[i]= q.Sql()
	}
	downScope := &Scope{}
	down(downScope)
	downScripts := make([]string, len(downScope.queries))
	for i, q := range downScope.queries {
		downScripts[i]= q.Sql()
	}
	m.Migrations = append(m.Migrations, Migration{
		Name: name,
		Version: version,
		UpScript: strings.Join(upScripts, ";"),
		DownScript: strings.Join(downScripts, ";"),
	})

	return nil
}

func (m *MigrateApplication) Run(args []string) {
	fmt.Println("Found ", len(m.Migrations), " migrations")
	fmt.Println(m.Migrations[0].UpScript)
	fmt.Println(m.Migrations[0].DownScript)
	fmt.Println("Running migrations")
}

func NewApp(dsn string) *MigrateApplication {
	sqlx.Connect("mysql", dsn)
	result := new(MigrateApplication)
	return result
}
