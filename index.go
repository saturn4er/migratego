package migrates

import (
	"strconv"
	"strings"
)

type IndexGenerator struct {
	name         string
	unique       bool
	columns      []*IndexColumn
	parser       string
	keyBlockSize int
	comment      string
}

func (i *IndexGenerator) Name(n string) *IndexGenerator {
	i.name = n
	return i
}
func (i *IndexGenerator) Columns(c ...*IndexColumn) *IndexGenerator {
	i.columns = append(i.columns, c...)
	return i
}
func (i *IndexGenerator) Comment(c string) *IndexGenerator {
	i.comment = c
	return i
}
func (i *IndexGenerator) KeyBlockSize(s int) *IndexGenerator {
	i.keyBlockSize = s
	return i
}
func (i *IndexGenerator) Parser(p string) *IndexGenerator {
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
	Column *ColumnGenerator
	Order  Order
	Length int
}

func (i *IndexColumn) Sql() string {
	var sql = wrapName(i.Column.name)
	if i.Length > 0 {
		sql += "(" + strconv.FormatInt(int64(i.Length), 10) + ")"
	}
	return sql + " " + string(i.Order)
}

// NewIndexColumn creates new IndexColumn
// Usage NewIndexColumn(column, orderType[optional], length[optional])
// orderType default value is ASC
// length default value is int
func NewIndexColumn(column *ColumnGenerator, params ...interface{}) *IndexColumn {
	var length int
	var order = OrderAsc
	var ok bool
	if len(params) > 0 {
		if order, ok = params[0].(Order); !ok {
			panic("first param should be of type `Order`")
		}
	}
	if len(params) > 1 {
		if length, ok = params[0].(int); !ok {
			panic("first param should be of type `int`")
		}
	}
	return &IndexColumn{
		Column: column,
		Order:  order,
		Length: length,
	}
}
func newIndexGenerator(name string, unique bool) *IndexGenerator {
	result := &IndexGenerator{
		name:   name,
		unique: unique,
	}
	return result
}
