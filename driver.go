package migratego

import (
	"github.com/saturn4er/migratego/databases/mysql"
	"github.com/saturn4er/migratego/types"
)

func shouldCheckDriver(driver string){
	if !checkDriver(driver){
		panic("We doesn't support "+driver+" driver")
	}
}
func checkDriver(driver string) bool {
	switch(driver){
	case "mysql":
		return true
	}
	return false
}
func getDriverQueryBuilder(driver string) QueryBuilder {
	switch driver {
	case "mysql":
		return new(mysql.MysqlQueryBuilder)
	}
	panic("Unknown driver:" + driver)
}
func getDriverClient(driver, dsn, transactionsTableName string) (types.DBClient, error) {
	switch driver {
	case "mysql":
		return mysql.NewClient(dsn, transactionsTableName)
	}
	panic("Unknown driver:" + driver)
}