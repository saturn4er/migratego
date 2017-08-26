package mysql

import "github.com/saturn4er/migratego"

func init() {
	migratego.DefineDriver("mysql", QueryBuilderConstructor, NewClient)
}

func QueryBuilderConstructor() QueryBuilder {
	return &MysqlQueryBuilder{}
}
