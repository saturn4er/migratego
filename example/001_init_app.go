package main

import (
	"github.com/saturn4er/migratego"
	_ "github.com/saturn4er/migratego/drivers/mysql"
)

func init() {
	app.AddMigration(1, "initApp", initAppUp, initAppDown)
}
func initAppUp(s migratego.QueryBuilder) {
	s.CreateTable("user", func(t migratego.CreateTableGenerator) {
		t.Column("id", "int").Primary()
		t.Column("name", "varchar(255)").NotNull()
		t.Column("password", "varchar(255)").NotNull()
		t.Charset("utf8mb4")
	})
	s.Table("user", func(scope migratego.TableScope) {
		scope.RemoveColumn("1")
	})
}
func initAppDown(s migratego.QueryBuilder) {
	s.DropTables("user").IfExists()
}
