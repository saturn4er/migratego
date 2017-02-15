package main

import (
	"os"

	"github.com/urfave/cli"
)

var Commands []cli.Command

const DefaultMigrationsFolder = "./migrations"

func main() {
	c := cli.NewApp()
	c.Commands = Commands
	c.Usage = "Tool, that helps developers to work with gomigrate"
	c.Version = "0.0.1"
	c.Run(os.Args)
}
