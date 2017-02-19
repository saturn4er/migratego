package migratego

import (
	"errors"
	"fmt"
	"sort"

	"github.com/saturn4er/migratego/types"
	"github.com/urfave/cli"
)

func RunToolCli(m *migrateApplication, args []string) error {
	tool := cli.NewApp()
	tool.HelpName = "migratego"
	client, err := m.getDriverClient()
	if err != nil {
		return err
	}
	err = client.PrepareTransactionsTable()
	if err != nil {
		return err
	}
	tool.Version = "1.0.0"
	tool.Usage = "Tool to manipulate with database versions"
	tool.Commands = []cli.Command{
		{
			Name:    "current",
			Aliases: []string{"c"},
			Usage:   "Migrations, that was applied to database",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "nwrap,nw",
					Usage: "Do not wrap the code",
				},
			},
			Action: func(c *cli.Context) error {
				applied, err := client.GetAppliedMigrations()
				if err != nil {
					return err
				}
				if len(applied) == 0 {
					fmt.Println("There's no migrations applied to database. Look's like it's empty")
					return nil
				}
				types.ShowMigrations(applied, c.Bool("nwrap"))
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "All available migrations",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "nwrap,nw",
					Usage: "Do not wrap the code",
				},
			},

			Action: func(c *cli.Context) error {
				applied, err := client.GetAppliedMigrations()
				if err != nil {
					return err
				}
				types.MergeMigrationsAppliedAt(m.migrations, applied)
				types.ShowMigrations(m.migrations,  c.Bool("nwrap"))
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
					Name:  "nwrap,nw",
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
					backupFilePath, err := client.Backup(backupPath)
					if err != nil {
						return err
					}
					fmt.Println("Backup created at ", backupFilePath)
				}
				sort.Sort(types.ByNumber(m.migrations))
				applied, err := client.GetAppliedMigrations()
				if err != nil {
					return err
				}

				toDowngrade, toUpgrade := types.FindWayBetweenMigrations(applied, m.migrations)
				if len(toDowngrade) == 0 && len(toUpgrade) == 0 {
					fmt.Println("Your database is already up-to-date")
					return nil
				}
				fmt.Println("Migrations, that will be applied:")
				types.ShowMigrationsToMigrate(toDowngrade, toUpgrade, c.Bool("nwrap"))

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
					err := client.ApplyMigration(&d, true)
					if err != nil {
						fmt.Println(err)
						return nil
					}
					err = client.RemoveMigration(&d)
					if err != nil {
						fmt.Println("Can't deletem migration from migrations table: ", err)
						return nil
					}
					fmt.Println("Ok!")
				}
				for _, d := range toUpgrade {
					fmt.Printf("Applying migration    %15d %15s....    ", d.Number, d.Name)
					err := client.ApplyMigration(&d, false)
					if err != nil {
						fmt.Println(err)
						return nil
					}
					err = client.InsertMigration(&d)
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
