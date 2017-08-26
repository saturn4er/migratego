package main

import (
	"github.com/saturn4er/migratego"
	_ "github.com/saturn4er/migratego/drivers/mysql"
)

func init() {
	app.AddMigration(1, "initApp", new(InitAppMigration))
}

type InitAppMigration struct{}

func (i *InitAppMigration) Up(qb migratego.QueryBuilder) {
	qb.CreateTable("user", func(t migratego.CreateTableGenerator) {
		t.Column("id", "int").Primary()
		t.Column("name", "varchar(255)").NotNull()
		t.Column("password", "varchar(255)").NotNull()
		t.Charset("utf8mb4")
	})
	qb.Table("user", func(scope migratego.TableScope) {
		scope.RemoveColumn("1")
	})
}
func (i *InitAppMigration) Down(qb migratego.QueryBuilder) {
	qb.DropTables("user").IfExists()
}
