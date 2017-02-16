package migrates

import (
	"fmt"
	"sort"

	"database/sql"

	"os"
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
			Usage:   "Show current version of database",
			Action: func(c *cli.Context) error {
				versions, err := mg.getCurrentVersion()
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("There's no migrations applied to database. Look's like it's empty")
						return nil
					}
					return err
				}

				fmt.Printf("Your databse is on %v #%v\n", versions.Name, versions.Number)
				return nil
			},
		},
		{
			Name:    "latest",
			Aliases: []string{"la"},
			Usage:   "Show latest migration",
			Action: func(c *cli.Context) error {
				sort.Sort(byVersion(m.migrations))
				if len(m.migrations) == 0 {
					fmt.Println("There's no migrations yet")
				}
				latestMigration := m.migrations[len(m.migrations)-1]
				fmt.Printf("Latest migration is %v #%v\n", latestMigration.Name, latestMigration.Version)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"li"},
			Usage:   "Show migrations list",
			Action: func(c *cli.Context) error {
				sort.Sort(byVersion(m.migrations))
				if len(m.migrations) == 0 {
					fmt.Println("There's no migrations yet")
				}
				appliedV, err := mg.getAllVersion()
				if err != nil {
				    return err
				}
				applieds := make(map[int]*Version)
				for _, a := range appliedV {
					ac := a
					applieds[a.Number] = &ac
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Version", "Name", "Applied date"})
				tableData := make([][]string, len(m.migrations))
				i := 0
				for _, m := range m.migrations {
					appliedDate := "---"
					if val, ok := applieds[m.Version]; ok{
						appliedDate = val.AppliedAt.Format("02-01-2016 15:04:05")
					}
					tableData[i] = []string{strconv.Itoa(m.Version), m.Name, appliedDate}
					i++
				}
				table.AppendBulk(tableData)
				table.Render()
				return nil
			},
		},
	}

	tool.Run(args)
	return nil
}
