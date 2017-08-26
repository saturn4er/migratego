package migratego

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func init() {
	command := cli.Command{
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
			client := c.App.Metadata["client"].(DBClient)
			applied, err := client.GetAppliedMigrations()
			if err != nil {
				return err
			}
			if len(applied) == 0 {
				fmt.Println("There's no migrations applied to database. Look's like it's empty")
				return nil
			}
			ShowMigrations(applied, c.Bool("nwrap"))
			return nil
		},
	}
	cliCommands = append(cliCommands, command)
}
