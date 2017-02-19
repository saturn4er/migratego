package main

import (
	"github.com/saturn4er/migratego"
	"github.com/saturn4er/migratego/types"
)


func init() {
	app.AddMigration(1, "initApp", initAppUp, initAppDown)
}
func initAppUp(s migratego.QueryBuilder) {
	s.CreateTable("user", func(t types.CreateTableGenerator) {
		t.Column("id", "int").Primary()
		t.Column("name", "varchar(255)").NotNull()
		t.Column("password", "varchar(255)").NotNull()
	})
}
func initAppDown(s migratego.QueryBuilder) {
	s.DropTables("user").IfExists()
}