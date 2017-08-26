package migratego

import (
	"gopkg.in/urfave/cli.v1"
)

func init() {
	command := cli.Command{
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
			client := c.App.Metadata["client"].(DBClient)
			app := c.App.Metadata["app"].(*migrateApplication)

			applied, err := client.GetAppliedMigrations()
			if err != nil {
				return err
			}
			MergeMigrationsAppliedAt(app.migrations, applied)
			ShowMigrations(app.migrations, c.Bool("nwrap"))
			return nil
		},
	}
	cliCommands = append(cliCommands, command)
}
