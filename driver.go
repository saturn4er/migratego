package migratego

type QueryBuilderConstructor func() QueryBuilder
type DBClientConstructor func(dsn, transactionTableName string) (DBClient, error)

type Driver struct {
	QueryBuilderConstructor QueryBuilderConstructor
	DBClientConstructor     DBClientConstructor
}

var drivers = make(map[string]Driver)

func DefineDriver(name string, qbc QueryBuilderConstructor, dbc DBClientConstructor) {
	if _, ok := drivers[name]; ok {
		panic("Driver '" + name + "' already defined")
	}
	drivers[name] = Driver{qbc, dbc}
}

func shouldCheckDriver(driver string) {
	if !checkDriver(driver) {
		panic("We doesn't support " + driver + " driver")
	}
}
func checkDriver(driver string) bool {
	_, ok := drivers[driver]
	return ok
}
func getDriverQueryBuilder(driver string) QueryBuilder {
	d, ok := drivers[driver]
	if !ok {
		panic("Unknown driver:" + driver)
	}
	return d.QueryBuilderConstructor()

}
func getDriverClient(driver, dsn, transactionsTableName string) (DBClient, error) {
	d, ok := drivers[driver]
	if !ok {
		panic("Unknown driver:" + driver)
	}
	return d.DBClientConstructor(dsn, transactionsTableName)
}
