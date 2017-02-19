package mysql

import (
	"strconv"
	"strings"
	"github.com/saturn4er/migratego/types"
)

type IndexGenerator struct {
	name         string
	unique       bool
	columns      []types.IndexColumnGenerator
	parser       string
	keyBlockSize int
	comment      string
}

func (i *IndexGenerator) Name(n string) types.IndexGenerator {
	i.name = n
	return i
}
func (i *IndexGenerator) Columns(c ...types.IndexColumnGenerator) types.IndexGenerator{
	i.columns = append(i.columns, c...)
	return i
}
func (i *IndexGenerator) Comment(c string) types.IndexGenerator {
	i.comment = c
	return i
}
func (i *IndexGenerator) KeyBlockSize(s int) types.IndexGenerator {
	i.keyBlockSize = s
	return i
}
func (i *IndexGenerator) Parser(p string) types.IndexGenerator {
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

type IndexColumn struct {
	Column types.ColumnGenerator
	Order  string
	Length int
}

func (i *IndexColumn) Sql() string {
	var sql = wrapName(i.Column.GetName())
	if i.Length > 0 {
		sql += "(" + strconv.FormatInt(int64(i.Length), 10) + ")"
	}
	return sql + " " + string(i.Order)
}

func newIndexGenerator(name string, unique bool) types.IndexGenerator {
	result := &IndexGenerator{
		name:   name,
		unique: unique,
	}
	return result
}
