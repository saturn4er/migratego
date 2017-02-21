package mysql

import (
	"strings"

	"github.com/saturn4er/migratego/types"
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
func (f *CreateTableColumn) AutoIncrement(primaryComment... string) types.CreateTableColumnGenerator {
	f.autoIncrement = true
	f.Primary(primaryComment...)
	return f
}

// Index add index to table for this column
func (f *CreateTableColumn) Index(name string, unique bool, order string, length int) types.IndexGenerator {
	index := newIndexGenerator(name, unique)

	index.Columns(&IndexColumnGenerator{
		Column: f.name,
		Order:  order,
		Length: length,
	})
	f.table.indexes = append(f.table.indexes, index)
	return index
}
func (f *CreateTableColumn) Primary(comment ...string) types.CreateTableColumnGenerator {
	var c string
	if len(comment) > 0 {
		c = comment[0]
	}
	f.table.primaryKey = NewPrimaryKeyGenerator([]string{f.name}, c)
	return f
}

// NotNull marks column as NOT NULL
func (f *CreateTableColumn) NotNull() types.CreateTableColumnGenerator {
	f.notNull = true
	return f
}

// Binary marks column as BINARY
func (f *CreateTableColumn) Binary() types.CreateTableColumnGenerator {
	f.binary = true
	return f
}
func (f *CreateTableColumn) ZeroFill() types.CreateTableColumnGenerator {
	f.zeroFill = true
	return f
}
func (f *CreateTableColumn) Unsigned() types.CreateTableColumnGenerator {
	f.unsigned = true
	return f
}
func (f *CreateTableColumn) Generated() types.CreateTableColumnGenerator {
	f.generated = true
	return f
}
func (f *CreateTableColumn) DefaultValue(v string) types.CreateTableColumnGenerator {
	f.defaultValue = v
	return f
}
func (f *CreateTableColumn) Comment(v string) types.CreateTableColumnGenerator {
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
	if f.defaultValue != "" {
		sql += " COMMENT '" + strings.Replace(f.comment, "'", "\\'", -1) + "'"
	}
	return sql
}
