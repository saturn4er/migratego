package mysql

import (
	"strings"

	"github.com/saturn4er/migratego/types"
)

type IndexType string

type Order string

type createTableGenerator struct {
	name          string
	engine        string
	comment       string
	charset       string
	columns       []types.ColumnGenerator
	indexes       []types.IndexGenerator
	primaryKey    *PrimaryKeyGenerator
	uniqueIndexes map[string]string
}

func (c *createTableGenerator) Sql() string {
	var sql = "CREATE TABLE `" + c.name + "`("

	w := make([]string, len(c.columns))
	for i, column := range c.columns {
		w[i] = column.Sql()
	}
	sql += strings.Join(w, ",")

	if c.primaryKey != nil {
		sql += ", " + c.primaryKey.Sql()
	}
	for _, index := range c.indexes {
		indexSql := index.Sql()
		if indexSql != "" {
			sql += ", " + indexSql
		}
	}

	sql += ")"
	if c.engine != "" {
		sql += " ENGINE = " + c.engine
	}
	if c.charset != "" {
		sql += " DEFAULT CHARACTER SET = " + c.charset
	}
	if c.engine != "" {
		sql += " COMMENT = '" + strings.Replace(c.comment, "'", "\\'", -1) + "'"
	}
	return sql
}
func (c *createTableGenerator) CharSet(charset string) types.CreateTableGenerator {
	c.charset = charset
	return c
}
func (c *createTableGenerator) Comment(comment string) types.CreateTableGenerator{
	c.comment = comment
	return c
}
func (c *createTableGenerator) Engine(engine string) types.CreateTableGenerator {
	c.engine = engine
	return c
}
func (c *createTableGenerator) Column(name string, Type string) types.ColumnGenerator {
	result := &ColumnGenerator{
		table: c,
		name:  name,
		fType: Type,
	}
	c.columns = append(c.columns, result)
	return result
}
func (c *createTableGenerator) Index(name string, unique bool) types.IndexGenerator {
	index := newIndexGenerator(name, unique)
	c.indexes = append(c.indexes, index)
	return index
}
// NewIndexColumn creates new IndexColumn
// Usage NewIndexColumn(column, order[optional], length[optional])
// orderType default value is ASC
// length default value is int
func (c *createTableGenerator) NewIndexColumn(column types.ColumnGenerator, params ...interface{}) types.IndexColumnGenerator {
	var length int
	var order = "ASC"
	var ok bool
	if len(params) > 0 {
		if order, ok = params[0].(string); !ok {
			panic("first param should be of type `string`")
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
func NewCreateTableGenerator(name string, sc func(types.CreateTableGenerator)) types.CreateTableGenerator {
	result := &createTableGenerator{
		name:          name,
		uniqueIndexes: make(map[string]string),
	}
	sc(result)
	return result
}
