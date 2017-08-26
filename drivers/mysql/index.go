package mysql

import (
	"strings"
)

type indexGenerator struct {
	name         string
	unique       bool
	columns      []IndexColumnGenerator
	parser       string
	keyBlockSize int
	comment      string
}

func (i *indexGenerator) Name(n string) IndexGenerator {
	i.name = n
	return i
}
func (i *indexGenerator) Columns(c ...IndexColumnGenerator) IndexGenerator {
	i.columns = append(i.columns, c...)
	return i
}
func (i *indexGenerator) Comment(c string) IndexGenerator {
	i.comment = c
	return i
}
func (i *indexGenerator) Unique() IndexGenerator {
	i.unique = true
	return i
}
func (i *indexGenerator) KeyBlockSize(s int) IndexGenerator {
	i.keyBlockSize = s
	return i
}
func (i *indexGenerator) Parser(p string) IndexGenerator {
	i.parser = p
	return i
}
func (i *indexGenerator) Sql() string {
	var sql string
	if len(i.columns) == 0 {
		return ""
	}
	if i.unique {
		sql += "UNIQUE "
	}
	columns := make([]string, len(i.columns))
	for i, c := range i.columns {
		columns[i] = c.Sql()
	}
	sql += "INDEX " + wrapName(i.name) + " (" + strings.Join(columns, ",") + ")"
	return sql
}

func newIndexGenerator(name string, unique bool) IndexGenerator {
	result := &indexGenerator{
		name:   name,
		unique: unique,
	}
	return result
}
