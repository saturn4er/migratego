package main

import "github.com/urfave/cli"

func init() {
	Commands = append(Commands, cli.Command{
		Name: "init",
		Usage: "Initialize new migrations project",
		Action: func(c *cli.Context) error{

			return nil
		},
	})
}
