package mysql

const (
	AlterTableAdd  = "ADD"
	AlterTableDrop = "DROP"
)

type AlterTableGenerator struct {
	table     string
	operation string
	query     Querier
}

func (a *AlterTableGenerator) Sql() string {
	sql := "ALTER TABLE " + wrapName(a.table)
	sql += " " + a.operation
	sql += " " + a.query.Sql()
	return sql
}
