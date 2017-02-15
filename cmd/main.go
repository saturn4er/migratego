package main

import (
	"os"

	"github.com/urfave/cli"
)

var Commands []cli.Command

func main() {
	c := cli.NewApp()
	c.Commands = Commands
	c.Run(os.Args)
}
