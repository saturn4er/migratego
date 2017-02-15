package migrates

import "strings"

type IndexType string

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc       = "DESC"
)

type CreateTable interface {
	Column(name string, Type string) *ColumnGenerator
	Index(name string, unique bool) *IndexGenerator
}
type createTable struct {
	name          string
	engine        string
	comment       string
	charset       string
	columns       []*ColumnGenerator
	indexes       []*IndexGenerator
	primaryKey    *PrimaryKeyGenerator
	uniqueIndexes map[string]string
}

func (c *createTable) Sql() string {
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
func (c *createTable) CharSet(charset string) CreateTable {
	c.charset = charset
	return c
}
func (c *createTable) Comment(comment string) CreateTable {
	c.comment = comment
	return c
}
func (c *createTable) Engine(engine string) CreateTable {
	c.engine = engine
	return c
}
func (c *createTable) Column(name string, Type string) *ColumnGenerator {
	result := &ColumnGenerator{
		table: c,
		name:  name,
		fType: Type,
	}
	c.columns = append(c.columns, result)
	return result
}
func (c *createTable) Index(name string, unique bool) *IndexGenerator {
	index := newIndexGenerator(name, unique)
	c.indexes = append(c.indexes, index)
	return index

}

func NewCreateTableGenerator(name string, sc func(scope CreateTable)) *createTable {
	result := &createTable{
		name:          name,
		uniqueIndexes: make(map[string]string),
	}
	sc(result)
	return result
}
