package mysql

import (
	"strings"

	"github.com/saturn4er/migratego/types"
)

type ColumnGenerators []ColumnGenerator

func (c ColumnGenerators) String() string {
	var result = make([]string, len(c))
	for i, v := range c {
		result[i] = "`" + strings.Replace(v.name, "`", "\\`", -1) + "`"
	}
	return strings.Join(result, ",")
}

type ColumnGenerator struct {
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
}

func (f *ColumnGenerator) GetName() string {
	return f.name
}

// AutoIncrement define column as AUTO_INCREMENT and new PRIMARY INDEX
func (f *ColumnGenerator) AutoIncrement(primaryComment string) types.ColumnGenerator {
	f.autoIncrement = true
	f.Primary(primaryComment)
	return f
}

// Index add index to table for this column
func (f *ColumnGenerator) Index(name string, unique bool, order string, length int) types.IndexGenerator {
	index := newIndexGenerator(name, unique)

	index.Columns(&IndexColumnGenerator{
		Column: f,
		Order:  order,
		Length: length,
	})
	f.table.indexes = append(f.table.indexes, index)
	return index
}
func (f *ColumnGenerator) Primary(comment ...string) types.ColumnGenerator {
	var c string
	if len(comment) > 0 {
		c = comment[0]
	}
	f.table.primaryKey = NewPrimaryKeyGenerator([]string{f.name}, c)
	return f
}

// NotNull marks column as NOT NULL
func (f *ColumnGenerator) NotNull() types.ColumnGenerator {
	f.notNull = true
	return f
}

// Binary marks column as BINARY
func (f *ColumnGenerator) Binary() types.ColumnGenerator {
	f.binary = true
	return f
}
func (f *ColumnGenerator) ZeroFill() types.ColumnGenerator {
	f.zeroFill = true
	return f
}
func (f *ColumnGenerator) Unsigned() types.ColumnGenerator {
	f.unsigned = true
	return f
}
func (f *ColumnGenerator) Generated() types.ColumnGenerator {
	f.generated = true
	return f
}
func (f *ColumnGenerator) DefaultValue(v string) types.ColumnGenerator {
	f.defaultValue = v
	return f
}
func (f *ColumnGenerator) Comment(v string) types.ColumnGenerator {
	f.comment = v
	return f
}

func (f *ColumnGenerator) Sql() string {
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
