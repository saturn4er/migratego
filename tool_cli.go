package migrates

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func RunToolCli(m *MigrateApplication, args []string) error {
	tool := cli.NewApp()
	mg, err := NewMigrator(m.dbVersionTable, m.dsn, m.migrations)
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
				current, err := mg.getCurrentVersion()
				if err != nil {
					return err
				}
				if current == nil {
					fmt.Println("There's no migrations applied to database. Look's like it's empty")
					return nil
				}

				fmt.Printf("Your databse is on %v #%v\n", current.Name, current.Number)
				return nil
			},
		},
		{
			Name:    "latest",
			Aliases: []string{"la"},
			Usage:   "Latest migration",
			Action: func(c *cli.Context) error {
				latestMigration := LatestMigration(m.migrations)
				if latestMigration == nil {
					fmt.Println("There's no migrations yet")
				}
				fmt.Printf("Latest migration is %v #%v\n", latestMigration.Name, latestMigration.Number)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"li"},
			Usage:   "Migrations list",
			Action: func(c *cli.Context) error {
				ShowMigrations(m.migrations, mg)
				return nil
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"mi"},
			Usage:   "Update database to actual version",
			Action: func(c *cli.Context) error {
				sort.Sort(byNumber(m.migrations))
				applied, err := mg.getAllVersions()
				if err != nil {
					return err
				}
				var toDowngrade []*Version
				var toUpgrade []*Migration
				var sameI = -1
				var tableData [][]string
				for i, a := range applied {
					if len(m.migrations)-1 < i || !a.SameAsMigration(&m.migrations[i]) {
						for j := len(applied) - 1; j >= i; j-- {
							apl := applied[j]
							appliedDate := apl.AppliedAt.Format("02-01-2016 15:04:05")
							toDowngrade = append(toDowngrade, &apl)
							tableData = append(tableData, []string{strconv.Itoa(apl.Number), apl.Name, "Down", apl.DownScript, appliedDate})
						}
						break
					}
					sameI = i
				}
				for i := sameI+1; i < len(m.migrations); i++ {
					mi := m.migrations[i]
					tableData = append(tableData, []string{strconv.Itoa(mi.Number), mi.Name, "Up", mi.UpScript, ""})
					toUpgrade = append(toUpgrade, &mi)
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"#", "Name", "Oper.", "Code", "Applied a" +
					"t"})
				table.SetColWidth(50)
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetRowLine(true)
				table.AppendBulk(tableData)
				table.Render()
				return nil
			},
		},
	}

	tool.Run(args)
	return nil
}
func LatestMigration(migrations []Migration) *Migration {
	if len(migrations) == 0 {
		return nil
	}
	sort.Sort(byNumber(migrations))
	return &migrations[len(migrations)-1]
}
func ShowMigrations(migrations []Migration, mg *Migrator) error {
	sort.Sort(byNumber(migrations))
	if len(migrations) == 0 {
		fmt.Println("There's no migrations yet")
	}
	appliedV, err := mg.getAllVersions()
	if err != nil {
		return err
	}
	applieds := make(map[int]*Version)
	for _, a := range appliedV {
		ac := a
		applieds[ac.Number] = &ac
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Name", "Applied at"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	tableData := make([][]string, len(migrations))
	i := 0
	for _, mi := range migrations {
		appliedDate := "---"
		if val, ok := applieds[mi.Number]; ok {

			appliedDate = val.AppliedAt.Format("02-01-2016 15:04:05")
		}
		tableData[i] = []string{strconv.Itoa(mi.Number), mi.Name, appliedDate}
		i++
	}
	table.AppendBulk(tableData)
	table.Render()
	return nil
}
