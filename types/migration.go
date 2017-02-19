package types

import "time"

type Migration struct {
	Number     int `db:"num"`
	Name       string
	UpScript   string     `db:"up_script"`
	DownScript string     `db:"down_script"`
	AppliedAt  *time.Time `db:"applied_at"`
}

func (v *Migration) Compare(m *Migration) bool {
	if v.Number != m.Number {
		return false
	}
	if v.Name != m.Name {
		return false
	}
	if v.UpScript != m.UpScript {
		return false
	}
	if v.DownScript != m.DownScript {
		return false
	}
	return true
}

type ByNumber []Migration

func (s ByNumber) Len() int {
	return len(s)
}

func (s ByNumber) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

func (s ByNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
