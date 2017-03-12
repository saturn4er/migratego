package mysql

import (
	"strings"

	"github.com/saturn4er/migratego"
)

type dropTablesGenerator struct {
	tables   []string
	ifExists bool
}

func (d *dropTablesGenerator) Table(tableName string) migratego.DropTablesGenerator {
	d.tables = append(d.tables, tableName)
	return d
}
func (d *dropTablesGenerator) IfExists() migratego.DropTablesGenerator {
	d.ifExists = true
	return d
}
func (d *dropTablesGenerator) Sql() string {
	if len(d.tables) == 0 {
		return ""
	}
	var sql = "DROP TABLE"
	if d.ifExists {
		sql += " IF EXISTS"
	}
	tableNames := make([]string, len(d.tables))
	for i, v := range d.tables {
		tableNames[i] = "`" + strings.Replace(v, "`", "\\`", -1) + "`"
	}
	sql += " " + strings.Join(tableNames, ",")
	return sql
}

func NewDropTablesGenerator(tableNames ...string) migratego.DropTablesGenerator {
	return &dropTablesGenerator{
		tables: tableNames,
	}
}
