package mysql

import "strconv"

type indexColumnGenerator struct {
	Column string
	Order  string
	Length int
}

func (i *indexColumnGenerator) Sql() string {
	var sql = wrapName(i.Column)
	if i.Length > 0 {
		sql += "(" + strconv.FormatInt(int64(i.Length), 10) + ")"
	}
	return sql + " " + string(i.Order)
}

func NewIndexColumnGenerator(column string, Order string, Length int) IndexColumnGenerator {
	return &indexColumnGenerator{
		Column: column,
		Order:  Order,
		Length: Length,
	}
}
