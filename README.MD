# MIGRATEGO

## Dependencies
- mysqldump

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


## Writing migrations

### Create table
```go
app.AddMigration(2, "initApp",
    func(s *migrates.Scope) {
        s.CreateTable("user", func(t migrates.CreateTable) {
            id := t.Column("id", "int").Primary()
            name := t.Column("name", "varchar(255)").NotNull()
            password := t.Column("password", "varchar(255)").NotNull()
        })
    },
    func(s *migrates.Scope) {
        s.DropTable("user").IfExists()
    },
)
```

## TODO

- AlterTable generator
- Add Postgres


