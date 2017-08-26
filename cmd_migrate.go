package migratego

import (
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"github.com/pkg/errors"
)

func init() {
	command := cli.Command{
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
			client := c.App.Metadata["client"].(DBClient)
			app := c.App.Metadata["app"].(*migrateApplication)

			backupPath := c.String("backup")
			if backupPath != "" {
				backupFilePath, err := client.Backup(backupPath)
				if err != nil {
					return err
				}
				fmt.Println("Backup created at ", backupFilePath)
			}
			SortMigrationsByNumber(app.migrations)
			applied, err := client.GetAppliedMigrations()
			if err != nil {
				return err
			}

			toDowngrade, toUpgrade := FindWayBetweenMigrations(applied, app.migrations)
			if len(toDowngrade) == 0 && len(toUpgrade) == 0 {
				fmt.Println("Your database is already up-to-date")
				return nil
			}
			fmt.Println("Migrations, that will be applied:")
			ShowMigrationsToMigrate(toDowngrade, toUpgrade, c.Bool("nwrap"))

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
	}
	cliCommands = append(cliCommands, command)
}
