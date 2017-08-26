package main

import (
	"github.com/saturn4er/migratego"
)

var app = migratego.NewApp("mysql", "root:password@tcp(127.0.0.1:3306)/catalogsrv")

func main() {
	app.Run()
}
