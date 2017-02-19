package migratego

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