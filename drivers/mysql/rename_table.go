package mysql

type RenameTableGenerator struct {
	oldName string
	newName string
}

func (s *RenameTableGenerator) Sql() string {
	return "RENAME TABLE " + wrapName(s.oldName) + " TO " + wrapName(s.newName)
}
