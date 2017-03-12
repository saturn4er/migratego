package mysql

type rawQuery string

func (s *rawQuery) Sql() string {
	return string(*s)
}
