package main

import (
	"fmt"

	"errors"

	"os"

	"path/filepath"
	"text/template"

	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize new migrations project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "dir,d",
				Usage: "Directory, where migrations will be created",
				Value: DefaultMigrationsFolder,
			},
			cli.StringFlag{
				Name:  "dsn,ds",
				Usage: "Data source name",
				Value: "root@127.0.0.1@/dbname",
			},
			cli.StringFlag{
				Name:  "table,t",
				Usage: "Table name, where migratego will store all info",
				Value: "",
			},
			cli.StringFlag{
				Name:  "driver,dr",
				Usage: "Database driver",
				Value: "mysql",
			},
		},
		Action: func(c *cli.Context) error {
			dir := c.String("dir")
			if dir == "" {
				fmt.Println("Please, specify dir")
				return nil
			}
			err := CreateMigrationsDirectory(dir)
			if err != nil {
				fmt.Println("Can't create migrations directory:", err)
				return nil
			}
			err = CreateDefaultMigrationsFiles(dir, c.String("driver"), c.String("dsn"), c.String("table"))
			if err != nil {
				fmt.Println("Can't create migrations files:", err)
				return nil
			}
			return nil
		},
	})
}

func CreateMigrationsDirectory(dir string) error {
	exists, err := dirExists(dir)
	if err != nil {
		return errors.New("can't check if path exists:" + err.Error())
	}
	if exists {
		empty, err := dirIsEmpty(dir)
		if err != nil {
			return errors.New("can't check if migrations directory is empty:" + err.Error())
		}
		if !empty {
			return errors.New(dir + " directory already exists and is not empty")
		}
	} else {
		err = os.MkdirAll(dir, 0764)
		if err != nil {
			return errors.New("can't create migrations directory: " + err.Error())
		}
		fmt.Println("Created " + dir)
	}

	return nil
}

func CreateDefaultMigrationsFiles(dir, driver, dsn, table string) error {
	mainTemplate := template.New("")
	_, err := mainTemplate.Parse(mainFileTpl)
	if err != nil {
		panic(err)
	}
	mainFilePath := filepath.Join(dir, "main.go")
	mainFile, err := os.OpenFile(mainFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0764)
	if err != nil {
		return errors.New("can't write to " + mainFilePath + ": " + err.Error())
	}
	err = mainTemplate.Execute(mainFile, map[string]string{
		"dsn":    dsn,
		"table":  table,
		"driver": driver,
	})
	if err != nil {
		return errors.New("can't write to " + mainFilePath + ": " + err.Error())
	}
	fmt.Println("Created " + mainFilePath)
	return nil
}
