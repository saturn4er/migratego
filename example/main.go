package main

import (
	"os"

	"github.com/saturn4er/migratego"
)

var app = migratego.NewApp("mysql", "root@tcp(192.168.99.100:3306)/dbname")
func main() {
	app.Run(os.Args)
}