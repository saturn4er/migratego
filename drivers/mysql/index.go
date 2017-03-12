package mysql

import (
	"strings"

	"github.com/saturn4er/migratego"
)

type IndexGenerator struct {
	name         string
	unique       bool
	columns      []migratego.IndexColumnGenerator
	parser       string
	keyBlockSize int
	comment      string
}

func (i *IndexGenerator) Name(n string) migratego.IndexGenerator {
	i.name = n
	return i
}
func (i *IndexGenerator) Columns(c ...migratego.IndexColumnGenerator) migratego.IndexGenerator {
	i.columns = append(i.columns, c...)
	return i
}
func (i *IndexGenerator) Comment(c string) migratego.IndexGenerator {
	i.comment = c
	return i
}
func (i *IndexGenerator) Unique() migratego.IndexGenerator {
	i.unique = true
	return i
}
func (i *IndexGenerator) KeyBlockSize(s int) migratego.IndexGenerator {
	i.keyBlockSize = s
	return i
}
func (i *IndexGenerator) Parser(p string) migratego.IndexGenerator {
	i.parser = p
	return i
}
func (i *IndexGenerator) Sql() string {
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

func newIndexGenerator(name string, unique bool) migratego.IndexGenerator {
	result := &IndexGenerator{
		name:   name,
		unique: unique,
	}
	return result
}
