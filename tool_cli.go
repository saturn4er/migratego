package migratego

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/saturn4er/migratego/types"
	"github.com/urfave/cli"
)

func RunToolCli(m *migrateApplication, args []string) error {
	tool := cli.NewApp()
	mg, err := newMigrator(m.dbVersionTable, m.driver, m.dsn, m.migrations)
	if err != nil {
		return err
	}
	err = mg.prepareDBVersionTable()
	if err != nil {
		return err
	}
	tool.Version = "1.0.0"
	tool.Usage = "Tool to manipulate with database versions"
	tool.Commands = []cli.Command{
		{
			Name:    "current",
			Aliases: []string{"c"},
			Usage:   "Current version of database",
			Action: func(c *cli.Context) error {
				applied, err := mg.getAppliedMigrations()
				if err != nil {
					return err
				}
				if len(applied) == 0 {
					fmt.Println("There's no migrations applied to database. Look's like it's empty")
					return nil
				}
				ShowMigrations(applied, true)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Migrations list",
			Action: func(c *cli.Context) error {
				applied, err := mg.getAppliedMigrations()
				if err != nil {
					return err
				}
				m := mergeMigrationsAppliedAt(m.migrations, applied)
				ShowMigrations(m, true)
				return nil
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Update database to actual version",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes,y",
					Usage: "Do not ask for confirmations",
				},
				cli.BoolFlag{
					Name:  "mwrap,nw",
					Usage: "Do not wrap the code",
				},
				cli.StringFlag{
					Name:  "backup,b",
					Usage: "Path to backup file",
				},
			},
			Action: func(c *cli.Context) error {
				backupPath := c.String("backup")
				if backupPath != "" {
					backupFilePath, err := mg.backupDatabase(backupPath)
					if err != nil {
						return err
					}
					fmt.Println("Backup created at ", backupFilePath)
				}
				sort.Sort(types.ByNumber(m.migrations))
				applied, err := mg.getAppliedMigrations()
				if err != nil {
					return err
				}

				toDowngrade, toUpgrade := findWayBetweenMigrations(applied, m.migrations)
				if len(toDowngrade) == 0 && len(toUpgrade) == 0 {
					fmt.Println("Your database is already up-to-date")
					return nil
				}
				fmt.Println("Migrations, that will be applied")
				ShowMigrationsToMigrate(toDowngrade, toUpgrade, c.Bool("mwrap"))

				if !c.Bool("y") {
					ok, err := askForConfirmation("Apply migrations?")
					if err != nil {
						return errors.New("can't obtain confirmation response")
					}
					if !ok {
						fmt.Println("Ok :(")
						return nil
					}
				}
				for _, d := range toDowngrade {
					fmt.Printf("Downgrading migration #%15d %15s....    ", d.Number, d.Name)
					err := mg.applyMigration(&d, true)
					if err != nil {
						fmt.Println(err)
						return nil
					}
					err = mg.deleteMigration(&d)
					if err != nil {
						fmt.Println("Can't deletem migration from migrations table: ", err)
						return nil
					}
					fmt.Println("Ok!")
				}
				for _, d := range toUpgrade {
					fmt.Printf("Applying migration    %15d %15s....    ", d.Number, d.Name)
					err := mg.applyMigration(&d, false)
					if err != nil {
						fmt.Println(err)
						return nil
					}
					err = mg.addMigration(&d)
					if err != nil {
						fmt.Println("Can't insert migration to migrations table: ", err)
						return nil
					}
					fmt.Println("Ok!")
				}
				return nil
			},
		},
	}

	tool.Run(args)
	return nil
}
func ShowMigrationsToMigrate(toDowngrade, toUpgrade []types.Migration, wrapCode bool) {
	var tableData [][]string
	for _, apl := range toDowngrade {
		applied := ""
		if apl.AppliedAt != nil {
			applied = apl.AppliedAt.Format("02-01-2016 15:04:05")
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
func ShowMigrations(migrations []types.Migration, showApplied bool) error {
	sort.Sort(types.ByNumber(migrations))
	if len(migrations) == 0 {
		fmt.Println("There's no migrations yet")
	}
	table := tablewriter.NewWriter(os.Stdout)
	var header = []string{"#", "Name", "Up", "Down"}
	if showApplied {
		header = append(header, "Applied at")
	}
	table.SetHeader(header)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetRowLine(true)
	tableData := make([][]string, len(migrations))
	i := 0
	for _, mi := range migrations {
		tableData[i] = []string{strconv.Itoa(mi.Number), mi.Name, mi.UpScript, mi.DownScript}
		if showApplied {
			var applied = ""
			if mi.AppliedAt != nil {
				applied = mi.AppliedAt.Format("02-01-2016 15:04:05")
			}
			tableData[i] = append(tableData[i], applied)
		}
		i++
	}
	table.AppendBulk(tableData)
	table.Render()
	return nil
}

// mergeMigrationsAppliedAt set applied at in
func mergeMigrationsAppliedAt(to []types.Migration, from []types.Migration) []types.Migration {
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
func findWayBetweenMigrations(applied, actual []types.Migration) (toDowngrade []types.Migration, toUpgrade []types.Migration) {
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
