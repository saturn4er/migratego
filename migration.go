package migratego

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

type DBMigration struct {
	Number     int `db:"num"`
	Name       string
	UpScript   string     `db:"up_script"`
	DownScript string     `db:"down_script"`
	AppliedAt  *time.Time `db:"applied_at"`
}

func (v *DBMigration) Compare(m *DBMigration) bool {
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

func SortMigrationsByNumber(migrations []DBMigration) {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number < migrations[j].Number
	})
}

func MergeMigrationsAppliedAt(to []DBMigration, from []DBMigration) []DBMigration {
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
func FindWayBetweenMigrations(applied, actual []DBMigration) (toDowngrade []DBMigration, toUpgrade []DBMigration) {
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
func ShowMigrationsToMigrate(toDowngrade, toUpgrade []DBMigration, wrapCode bool) {
	var tableData [][]string
	for _, apl := range toDowngrade {
		applied := ""
		if apl.AppliedAt != nil {
			applied = apl.AppliedAt.Format("02.01.06 15:04:05")
		}
		code := apl.DownScript
		if len(code) > 47 && !wrapCode {
			code = code[:47] + "..."
		}
		tableData = append(tableData, []string{strconv.Itoa(apl.Number), apl.Name, "Down", code, applied})
	}
	for _, apl := range toUpgrade {
		code := apl.UpScript
		if len(code) > 47 && !wrapCode {
			code = code[:47] + "..."
		}
		tableData = append(tableData, []string{strconv.Itoa(apl.Number), apl.Name, "Up", code, ""})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Name", "Oper.", "Code", "Applied at"})
	table.SetColWidth(50)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetRowLine(true)
	table.AppendBulk(tableData)
	table.Render()
}
func ShowMigrations(migrations []DBMigration, wrap bool) error {
	if len(migrations) == 0 {
		fmt.Println("There's no migrations yet")
	}
	SortMigrationsByNumber(migrations)
	table := tablewriter.NewWriter(os.Stdout)
	var header = []string{"#", "Name", "Up", "Down", "Applied at"}
	table.SetHeader(header)
	table.SetColWidth(50)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetRowLine(true)
	tableData := make([][]string, len(migrations))
	for _, mi := range migrations {
		var applied = ""
		if mi.AppliedAt != nil {
			applied = mi.AppliedAt.Format("02-01-2016 15:04:05")
		}
		tableData = append(tableData, []string{strconv.Itoa(mi.Number), mi.Name, wrapCode(mi.UpScript, wrap), wrapCode(mi.DownScript, wrap), applied})
	}
	table.AppendBulk(tableData)
	table.Render()
	return nil
}
func wrapCode(code string, wrap bool) string {
	if len(code) > 47 && !wrap {
		code = code[:47] + "..."
	}
	return code
}
