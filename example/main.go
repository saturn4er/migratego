package main

import (
	"os"

	"github.com/saturn4er/migratego"
)

func main() {
	app := migrates.NewApp("root@127.0.0.1@/dbname")
	app.AddMigration(1, "initApp",
		func(s *migrates.Scope) {
			s.CreateTable("hello", func(t migrates.CreateTable) {
				id := t.Column("id", "int").Primary()
				g := t.Column("g", "varchar(255)").NotNull()
				t.Index("123", false).Columns(
					migrates.NewIndexColumn(g),
					migrates.NewIndexColumn(id),
				)
			})
		},
		func(s *migrates.Scope) {
			s.DropTable("hello").IfExists()
		},
	)
	app.Run(os.Args)
}
