package mysql

import (
	"strings"

	"github.com/saturn4er/migratego"
)

type UpdateTableAddColumnGenerator struct {
	tableScope    *TableScope
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

func (f *UpdateTableAddColumnGenerator) GetName() string {
	return f.name
}

func (f *UpdateTableAddColumnGenerator) Index(name string, unique bool, order string, length int) migratego.IndexGenerator {
	index := f.tableScope.AddIndex(name, unique)
	index.Columns(&IndexColumnGenerator{
		Column: f.name,
		Order:  order,
		Length: length,
	})
	return index
}

// NotNull marks column as NOT NULL
func (f *UpdateTableAddColumnGenerator) NotNull() migratego.UpdateTableAddColumnGenerator {
	f.notNull = true
	return f
}

// Binary marks column as BINARY
func (f *UpdateTableAddColumnGenerator) Binary() migratego.UpdateTableAddColumnGenerator {
	f.binary = true
	return f
}
func (f *UpdateTableAddColumnGenerator) ZeroFill() migratego.UpdateTableAddColumnGenerator {
	f.zeroFill = true
	return f
}
func (f *UpdateTableAddColumnGenerator) Unsigned() migratego.UpdateTableAddColumnGenerator {
	f.unsigned = true
	return f
}
func (f *UpdateTableAddColumnGenerator) Generated() migratego.UpdateTableAddColumnGenerator {
	f.generated = true
	return f
}
func (f *UpdateTableAddColumnGenerator) DefaultValue(v string) migratego.UpdateTableAddColumnGenerator {
	f.defaultValue = v
	return f
}
func (f *UpdateTableAddColumnGenerator) Comment(v string) migratego.UpdateTableAddColumnGenerator {
	f.comment = v
	return f
}
func (f *UpdateTableAddColumnGenerator) After(field string) migratego.UpdateTableAddColumnGenerator {
	f.after = field
	return f
}

func (f *UpdateTableAddColumnGenerator) Sql() string {
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
