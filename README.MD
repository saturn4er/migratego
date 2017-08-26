[![Build Status](https://travis-ci.org/saturn4er/migratego.svg?branch=master)](https://travis-ci.org/saturn4er/migratego)
[![Coverage Status](https://coveralls.io/repos/github/saturn4er/migratego/badge.svg?branch=master)](https://coveralls.io/github/saturn4er/migratego?branch=master)
# MIGRATEGO

Library to write and build migrations using Golang.

## Dependencies
- Go 1.9+
- mysqldump ( to dump mysql database before applying migrations )

## Installation

**Library**

    go get github.com/saturn4er/migratego

**Cli tool**

    go get github.com/saturn4er/migratego/cmd/migratego

## Usage

* Init your migration application with `migratego init`
* Create your first migration with `migratego g initial_migration`
* Open `./migrations/001_initial_migration.go` file and write your first migration
* Build it with `go build -o migrations ./migrations/*.go`
* Apply your first `initial_migration` to database with `./migrations m`


## Examples

### Create table
```go
package main

import (
	"github.com/saturn4er/migratego"
	_ "github.com/saturn4er/migratego/drivers/mysql"
)

func init() {
	app.AddMigration(1, "initApp", new(InitAppMigration))
}

type InitAppMigration struct{}

func (i *InitAppMigration) Up(qb migratego.QueryBuilder) {
	qb.CreateTable("user", func(t migratego.CreateTableGenerator) {
		t.Column("id", "int").Primary()
		t.Column("name", "varchar(255)").NotNull()
		t.Column("password", "varchar(255)").NotNull()
		t.Charset("utf8mb4")
	})
	qb.Table("user", func(scope migratego.TableScope) {
		scope.RemoveColumn("1")
	})
}
func (i *InitAppMigration) Down(qb migratego.QueryBuilder) {
	qb.DropTables("user").IfExists()
}
```