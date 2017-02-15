package migrates

import (
	"errors"

	"fmt"
	"sort"

	"github.com/urfave/cli"
)

func init() {
	var upCommand = cli.Command{
		Name:  "up",
		Usage: "Upgrade database version to next version",
	}
	var downCommand = cli.Command{
		Name:  "down",
		Usage: "Downgrade database version to previous version",
	}
	var migrateCommand = cli.Command{
		Name:        "migrate",
		Aliases:     []string{"m"},
		Usage:       "Update database to the latest version",
		Subcommands: []cli.Command{upCommand, downCommand},
		Action: func(c *cli.Context) error {
			migrations, ok := c.App.Metadata["migrations"].([]Migration)
			if !ok || len(migrations)==0 {
				return errors.New("can't find any migrations")
			}
			sort.Sort(byVersion(migrations))
			fmt.Printf("Found %d migrations\n", len(migrations))
			latestMigration := migrations[len(migrations)-1]
			fmt.Printf("Latest migration is %s(#%d)\n",latestMigration.Name, latestMigration.Version)
			return errors.New("123")
		},
	}

	commands = append(commands, migrateCommand)
}

func MigrateToVersion(version int) error {
	fmt.Println(version)
	return nil
}
