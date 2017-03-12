package mysql

import (
	"strings"

	"github.com/saturn4er/migratego"
)

type CreateTableColumns []CreateTableColumn

func (c CreateTableColumns) String() string {
	var result = make([]string, len(c))
	for i, v := range c {
		result[i] = "`" + strings.Replace(v.name, "`", "\\`", -1) + "`"
	}
	return strings.Join(result, ",")
}

type CreateTableColumn struct {
	table *createTableGenerator
	name  string
	fType string

	binary        bool
	unsigned      bool
	zeroFill      bool
	autoIncrement bool
	notNull       bool
	generated     bool
	comment       string
	defaultValue  string
	charset       string
	after         string
}

func (f *CreateTableColumn) GetName() string {
	return f.name
}

// AutoIncrement define column as AUTO_INCREMENT and new PRIMARY INDEX
func (f *CreateTableColumn) AutoIncrement(primaryComment ...string) migratego.CreateTableColumnGenerator {
	f.autoIncrement = true
	f.Primary(primaryComment...)
	return f
}

// Index add index to table for this column
// Usage: Index("index_name", true, "DESC", 10)
// This will create unique index "index_name" and it will add this column to it
func (f *CreateTableColumn) Index(name string, unique bool, params ...interface{}) migratego.IndexGenerator {
	if name == "" {
		name = "idx_" + f.name
	}
	index := newIndexGenerator(name, unique)
	var order = "ASC"
	var length = 0
	if len(params) > 0 {
		o, ok := params[0].(string)
		if !ok {
			panic("Third param should be string (Order)")
		}
		if o != "" {
			order = o
		}
	}
	if len(params) > 1 {
		var ok bool
		length, ok = params[0].(int)
		if !ok {
			panic("Fourth param should be int (Length)")
		}

	}

	index.Columns(&IndexColumnGenerator{
		Column: f.name,
		Order:  order,
		Length: length,
	})
	f.table.indexes = append(f.table.indexes, index)
	return index
}
func (f *CreateTableColumn) Primary(comment ...string) migratego.CreateTableColumnGenerator {
	var c string
	if len(comment) > 0 {
		c = comment[0]
	}
	f.table.primaryKey = NewPrimaryKeyGenerator([]string{f.name}, c)
	return f
}

// NotNull marks column as NOT NULL
func (f *CreateTableColumn) NotNull() migratego.CreateTableColumnGenerator {
	f.notNull = true
	return f
}

// Binary marks column as BINARY
func (f *CreateTableColumn) Binary() migratego.CreateTableColumnGenerator {
	f.binary = true
	return f
}
func (f *CreateTableColumn) ZeroFill() migratego.CreateTableColumnGenerator {
	f.zeroFill = true
	return f
}
func (f *CreateTableColumn) Unsigned() migratego.CreateTableColumnGenerator {
	f.unsigned = true
	return f
}
func (f *CreateTableColumn) Generated() migratego.CreateTableColumnGenerator {
	f.generated = true
	return f
}
func (f *CreateTableColumn) DefaultValue(v string) migratego.CreateTableColumnGenerator {
	f.defaultValue = v
	return f
}
func (f *CreateTableColumn) Comment(v string) migratego.CreateTableColumnGenerator {
	f.comment = v
	return f
}

func (f *CreateTableColumn) Sql() string {
	sql := "`" + f.name + "` " + string(f.fType)
	if f.unsigned {
		sql += " UNSIGNED"
	}
	if f.zeroFill {
		sql += " ZEROFILL"
	}
	if f.binary {
		sql += " BINARY"
	}
	if f.notNull {
		sql += " NOT NULL"
	} else {
		sql += " NULL"
	}
	if f.defaultValue != "" {
		if f.generated {
			sql += " GENERATED ALWAYS AS ('" + strings.Replace(f.defaultValue, "'", "\\'", -1) + "')"
		} else {
			sql += " DEFAULT '" + strings.Replace(f.defaultValue, "'", "\\'", -1) + "'"
		}
	}
	if f.autoIncrement {
		sql += " AUTO_INCREMENT"
	}
	if f.charset != "" {
		sql += "CHARACTER SET '" + string(f.charset) + "' NULL"
	}
	if f.comment != "" {
		sql += " COMMENT '" + strings.Replace(f.comment, "'", "\\'", -1) + "'"
	}
	return sql
}
