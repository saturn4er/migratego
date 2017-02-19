package mysql

import (
	"strconv"

	"github.com/saturn4er/migratego/types"
)

type IndexColumnGenerator struct {
	Column types.ColumnGenerator
	Order  string
	Length int
}

func (i *IndexColumnGenerator) Sql() string {
	var sql = wrapName(i.Column.GetName())
	if i.Length > 0 {
		sql += "(" + strconv.FormatInt(int64(i.Length), 10) + ")"
	}
	return sql + " " + string(i.Order)
}

func NewIndexColumnGenerator(column types.ColumnGenerator, Order string, Length int) types.IndexColumnGenerator {
	return &IndexColumnGenerator{
		Column: column,
		Order:  Order,
		Length: Length,
	}
}
