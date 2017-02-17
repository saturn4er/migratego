package migrates

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

type byNumber []Migration

func (s byNumber) Len() int {
	return len(s)
}

func (s byNumber) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

func (s byNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
