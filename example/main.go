package main

import (
	"os"

	"github.com/saturn4er/migratego"
)

var app = migratego.NewApp("mysql", "root@tcp(127.0.0.1:3306)/dbname")

func main() {
	app.Run(os.Args)
}
