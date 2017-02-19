package mysql

import "strings"

type PrimaryKeyGenerator struct {
	columns []string
	comment string
}

func (p *PrimaryKeyGenerator) Sql() string {
	var sql = "PRIMARY KEY (" + wrapNames(p.columns) + ")"
	if p.comment != "" {
		sql += " COMMENT '" + strings.Replace(p.comment, "'", "\\'", -1) + "'"
	}
	return sql
}

func NewPrimaryKeyGenerator(columns []string, comment ...string) *PrimaryKeyGenerator {
	result := new(PrimaryKeyGenerator)
	result.columns = columns
	if len(comment) != 0 {
		result.comment = comment[0]
	}
	return result
}
