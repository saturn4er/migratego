package mysql

import (
	"strconv"

	"github.com/saturn4er/migratego"
)

type IndexColumnGenerator struct {
	Column string
	Order  string
	Length int
}

func (i *IndexColumnGenerator) Sql() string {
	var sql = wrapName(i.Column)
	if i.Length > 0 {
		sql += "(" + strconv.FormatInt(int64(i.Length), 10) + ")"
	}
	return sql + " " + string(i.Order)
}

func NewIndexColumnGenerator(column string, Order string, Length int) migratego.IndexColumnGenerator {
	return &IndexColumnGenerator{
		Column: column,
		Order:  Order,
		Length: Length,
	}
}
