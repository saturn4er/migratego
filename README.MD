[![Build Status](https://travis-ci.org/saturn4er/migratego.svg?branch=master)](https://travis-ci.org/saturn4er/migratego)
[![Coverage Status](https://coveralls.io/repos/github/saturn4er/migratego/badge.svg?branch=master)](https://coveralls.io/github/saturn4er/migratego?branch=master)
# MIGRATEGO

Library to write and build migrations using Golang.

## Dependencies
- Go 1.7+
- mysqldump ( to dump mysql database before applying migrations

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
	"github.com/saturn4er/migratego/types"
)


func init() {
	app.AddMigration(2, "initApp", initAppUp, initAppDown)
}
func initAppUp(s migratego.QueryBuilder) {
	s.CreateTable("user", func(t types.CreateTableGenerator) {
		t.Column("id", "int").Primary()
		t.Column("name", "varchar(255)").NotNull()
		t.Column("password", "varchar(255)").NotNull()

	})
}
func initAppDown(s migratego.QueryBuilder) {
	s.DropTables("user").IfExists()
}
```

## Sql generators
 * Drop tables
 * Create Table
   * Column
     * Primary
     * Not null
     * Unsigned
     * Binary
     * Zero Fill
     * Generated
     * Default Value
     * Comment
     * Index
     * Auto Increment
   * Index
   * Foreign key (TODO)
 * Table
   * Rename
   * Delete
   * Add column
   * Add index
   * Remove index
   * Update column (TODO)
 * Raw Query
 

## TODO
- Add Postgres
- Foreign Key
- Tool to Create initial migration by existing tables