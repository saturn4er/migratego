package main

import (
	"fmt"

	"path/filepath"
	"strconv"

	"os"

	"errors"

	"text/template"

	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "generate",
		Usage:     "Generates migration",
		ArgsUsage: "[migration_version int] [migration_name string]",
		Aliases:   []string{"g"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "dir,d",
				Usage: "Directory, where migrations will be created",
				Value: DefaultMigrationsFolder,
			},
		},
		Action: func(c *cli.Context) error {
			args := c.Args()
			if len(args) < 2 {
				cli.ShowSubcommandHelp(c)
				return nil
			}
			dir := c.String("dir")
			exists, err := dirExists(dir)
			if err != nil {
				fmt.Println("Can't check, if migrations directory exists:", err.Error())
				return nil
			}
			if !exists {
				fmt.Println("Migrations directory(" + dir + ") doesn't exist. Please, run `" + c.App.Name + " init`")
				return nil
			}
			version, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Bad version(" + args[0] + "). It should be of integer type ")
				return nil
			}

			var migrationArgs []string
			if len(args) > 2 {
				migrationArgs = args[2:]
			}
			return generateMigration(version, args[1], dir, migrationArgs)
		},
	})
}

func generateMigration(version int, name string, dir string, args []string) error {
	migrationFilePath := filepath.Join(dir, fmt.Sprintf("%03d_%s.go", version, name))

	exists, err := fileExists(migrationFilePath)
	if err != nil {
		return errors.New("can't check, if migration file already exists:" + err.Error())
	}
	if exists {
		return errors.New("migration file (" + migrationFilePath + ") already exists")
	}
	migrationFile, err := os.OpenFile(migrationFilePath, os.O_CREATE|os.O_WRONLY, 0764)
	if err != nil {
		return errors.New("Can't create migration file: " + err.Error())
	}

	mainTemplate, err := template.New("").Parse(migrationFileTemplate)
	if err != nil {
		panic(err)
	}
	upBody, upImports := GetUpMigrationBodyByName(name)
	downBody, downImports := GetDownMigrationBodyByName(name)
	usedImports := make(map[string]bool)
	usedImports["github.com/saturn4er/migratego"] = true
	var imports []string
	for _, i := range upImports {
		if !usedImports[i] {
			imports = append(imports, i)
			usedImports[i] = true
		}
	}
	for _, i := range downImports {
		if !usedImports[i] {
			imports = append(imports, i)
			usedImports[i] = true
		}
	}
	err = mainTemplate.Execute(migrationFile, map[string]interface{}{
		"version":  version,
		"name":     name,
		"upBody":   upBody,
		"downBody": downBody,
		"imports":  imports,
	})
	if err != nil {
		return errors.New("can't write to " + migrationFilePath + ": " + err.Error())
	}
	return nil
}

// GetUpMigrationBodyByName TODO
func GetUpMigrationBodyByName(name string) (string, []string) {
	return "", nil
}

// GetUpMigrationBodyByName TODO
func GetDownMigrationBodyByName(name string) (string, []string) {
	return "", nil
}
