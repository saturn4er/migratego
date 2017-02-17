package migrates

import "strings"

type PrimaryKeyGenerator struct {
	indexName string
	columns   []string
	comment   string
}

func (p *PrimaryKeyGenerator) Sql() string {
	var sql = "PRIMARY KEY (" + wrapNames(p.columns) + ")"
	if p.comment != "" {
		sql += " COMMENT '" + strings.Replace(p.comment, "'", "\\'", -1) + "'"
	}
	return sql
}

func NewPrimaryKeyGenerator(columns []string, comment string) *PrimaryKeyGenerator {
	return &PrimaryKeyGenerator{
		columns: columns,
		comment: comment,
	}
}
