package mysql

import "github.com/saturn4er/migratego"

const (
	AlterTableAdd  = "ADD"
	AlterTableDrop = "DROP"
)

type AlterTableGenerator struct {
	table     string
	operation string
	query     migratego.Querier
}

func (a *AlterTableGenerator) Sql() string {
	sql := "ALTER TABLE " + wrapName(a.table)
	sql += " " + a.operation
	sql += " " + a.query.Sql()
	return sql
}
