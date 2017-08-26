package mysql

type tableScope struct {
	name    string
	builder *MysqlQueryBuilder
}
type UpdateTableIndexes struct {
}

func (t *tableScope) AddColumn(name string, Type string) UpdateTableAddColumnGenerator {
	if name == "" {
		panic("Can't add column to table with empty name")
	}
	cg := &updateTableAddColumnGenerator{
		tableScope: t,
		name:       name,
		fType:      Type,
	}
	t.builder.generators = append(t.builder.generators, cg)
	return cg
}
func (t *tableScope) RemoveColumn(name string) TableScope {
	q := rawQuery("COLUMN " + wrapName(name))
	g := AlterTableGenerator{
		table:     t.name,
		operation: AlterTableDrop,
		query:     &q,
	}
	t.builder.generators = append(t.builder.generators, &g)
	return t
}

func (t *tableScope) Rename(newName string) TableScope {
	if newName == "" {
		panic("New name of table should not be empty")
	}
	t.builder.generators = append(t.builder.generators, &RenameTableGenerator{oldName: t.name, newName: newName})
	t.name = newName
	return t
}
func (t *tableScope) AddIndex(name string, unique bool) IndexGenerator {
	index := newIndexGenerator(name, unique)

	g := AlterTableGenerator{
		table:     t.name,
		operation: AlterTableAdd,
		query:     index,
	}
	t.builder.generators = append(t.builder.generators, &g)
	return index
}
func (t *tableScope) RemoveIndex(name string) TableScope {
	q := rawQuery("INDEX " + wrapName(name))
	g := AlterTableGenerator{
		table:     t.name,
		operation: AlterTableAdd,
		query:     &q,
	}
	t.builder.generators = append(t.builder.generators, &g)
	return t
}
func (t *tableScope) Delete() {
	t.builder.generators = append(t.builder.generators, &dropTablesGenerator{tables: []string{t.name}})
}
func (t *tableScope) DeleteIfExists() {
	t.builder.generators = append(t.builder.generators, &dropTablesGenerator{ifExists: true, tables: []string{t.name}})
}
