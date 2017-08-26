package migratego

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func askForConfirmation(s string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}

func ShowMigrations(migrations []DBMigration, wrap bool) error {
	if len(migrations) == 0 {
		return errors.New("No migrations :(")
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
