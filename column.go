package migrates

import "strings"

type ColumnGenerators []ColumnGenerator

func (c ColumnGenerators) String() string {
	var result = make([]string, len(c))
	for i, v := range c {
		result[i] = "`" + strings.Replace(v.name, "`", "\\`", -1) + "`"
	}
	return strings.Join(result, ",")
}

type ColumnGenerator struct {
	table *createTable
	name  string
	fType string

	unique        bool
	binary        bool
	unsigned      bool
	zeroFill      bool
	autoIncrement bool
	notNull       bool
	generated     bool
	comment       *string
	defaultValue  *string
	charset       *string
}

func (f *ColumnGenerator) Unique() *ColumnGenerator {
	f.unique = true
	return f
}

// AutoIncrement define column as AUTO_INCREMENT and new PRIMARY INDEX
func (f *ColumnGenerator) AutoIncrement(primaryComment string) *ColumnGenerator {
	f.autoIncrement = true
	f.Primary(primaryComment)
	return f
}

// Index add index to table for this column
func (f *ColumnGenerator) Index(name string, unique bool, order Order, length int) *IndexGenerator {
	index := newIndexGenerator(name, unique)
	index.Columns(NewIndexColumn(f, order, length))
	f.table.indexes = append(f.table.indexes, index)
	return index
}
func (f *ColumnGenerator) Primary(comment ...string) *ColumnGenerator {
	var c string
	if len(comment)>0{
		c = comment[0]
	}
	f.table.primaryKey = NewPrimaryKeyGenerator([]string{f.name}, c)
	return f
}

// NotNull marks column as NOT NULL
func (f *ColumnGenerator) NotNull() *ColumnGenerator {
	f.notNull = true
	return f
}

// Binary marks column as BINARY
func (f *ColumnGenerator) Binary() *ColumnGenerator {
	f.binary = true
	return f
}
func (f *ColumnGenerator) ZeroFill() *ColumnGenerator {
	f.zeroFill = true
	return f
}
func (f *ColumnGenerator) Unsigned() *ColumnGenerator {
	f.unsigned = true
	return f
}
func (f *ColumnGenerator) Generated() *ColumnGenerator {
	f.generated = true
	return f
}
func (f *ColumnGenerator) DefaultValue(v string) *ColumnGenerator {
	f.defaultValue = &v
	return f
}
func (f *ColumnGenerator) Comment(v string) *ColumnGenerator {
	f.comment = &v
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
	if f.defaultValue != nil {
		if f.generated {
			sql += " GENERATED ALWAYS AS ('" + strings.Replace(*f.defaultValue, "'", "\\'", -1) + "')"
		} else {
			sql += " DEFAULT '" + strings.Replace(*f.defaultValue, "'", "\\'", -1) + "'"
		}
	}
	if f.autoIncrement {
		sql += " AUTO_INCREMENT"
	}
	if f.charset != nil {
		sql += "CHARACTER SET '" + string(*f.charset) + "' NULL"
	}
	if f.defaultValue != nil {
		sql += " COMMENT '" + strings.Replace(*f.comment, "'", "\\'", -1) + "'"
	}
	return sql
}
