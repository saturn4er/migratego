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
func MergeMigrationsAppliedAt(to []Migration, from []Migration) []Migration {
	for i, mTo := range to {
		for _, mFrom := range from {
			if mTo.Compare(&mFrom) {
				to[i].AppliedAt = mFrom.AppliedAt
				break
			}
		}
	}
	return to
}
// findWayBetweenMigrations find path between two migrations list.
func FindWayBetweenMigrations(applied, actual []Migration) (toDowngrade []Migration, toUpgrade []Migration) {
	var sameI = -1
	for i, a := range applied {
		if len(actual)-1 < i || !a.Compare(&actual[i]) {
			for j := len(applied) - 1; j >= i; j-- {
				toDowngrade = append(toDowngrade, applied[j])
			}
			break
		}
		sameI = i
	}
	for i := sameI + 1; i < len(actual); i++ {
		toUpgrade = append(toUpgrade, actual[i])
	}
	return
}