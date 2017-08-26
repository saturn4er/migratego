package mysql

import (
	"strings"
)

type updateTableAddColumnGenerator struct {
	tableScope    *tableScope
	name          string
	fType         string
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

func (f *updateTableAddColumnGenerator) GetName() string {
	return f.name
}

func (f *updateTableAddColumnGenerator) Index(name string, unique bool, order string, length int) IndexGenerator {
	index := f.tableScope.AddIndex(name, unique)
	index.Columns(&indexColumnGenerator{
		Column: f.name,
		Order:  order,
		Length: length,
	})
	return index
}

// NotNull marks column as NOT NULL
func (f *updateTableAddColumnGenerator) NotNull() UpdateTableAddColumnGenerator {
	f.notNull = true
	return f
}

// Binary marks column as BINARY
func (f *updateTableAddColumnGenerator) Binary() UpdateTableAddColumnGenerator {
	f.binary = true
	return f
}
func (f *updateTableAddColumnGenerator) ZeroFill() UpdateTableAddColumnGenerator {
	f.zeroFill = true
	return f
}
func (f *updateTableAddColumnGenerator) Unsigned() UpdateTableAddColumnGenerator {
	f.unsigned = true
	return f
}
func (f *updateTableAddColumnGenerator) Generated() UpdateTableAddColumnGenerator {
	f.generated = true
	return f
}
func (f *updateTableAddColumnGenerator) DefaultValue(v string) UpdateTableAddColumnGenerator {
	f.defaultValue = v
	return f
}
func (f *updateTableAddColumnGenerator) Comment(v string) UpdateTableAddColumnGenerator {
	f.comment = v
	return f
}
func (f *updateTableAddColumnGenerator) After(field string) UpdateTableAddColumnGenerator {
	f.after = field
	return f
}

func (f *updateTableAddColumnGenerator) Sql() string {
	sql := "ALTER TABLE " + wrapName(f.tableScope.name) + " ADD COLUMN " + wrapName(f.name) + " " + string(f.fType)
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
	if f.after != "" {
		sql += " AFTER " + wrapName(f.after)
	}
	return sql
}
