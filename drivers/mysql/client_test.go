package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMysqlClient(t *testing.T) {
	Convey("MysqlDBClient", t, func() {
		d, err := NewClient("root:password@tcp(127.0.0.1:3306)/migratego_test", "schema_version")
		So(err, ShouldBeNil)
		So(d, ShouldNotBeNil)
		c := d.(*MysqlClient)
		Convey("Migrations table shouldn't exists", func() {
			exists, err := c.dbVersionTableExists()
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)

		})
		Convey("Should create migrations table", func() {
			So(c.PrepareTransactionsTable(), ShouldBeNil)
			exists, err := c.dbVersionTableExists()
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
		})
		testMigration := &DBMigration{
			Name:   "test migration",
			Number: 1,
			UpScript: NewCreateTableGenerator("test_table", func(c CreateTableGenerator) {
				c.Column("id", "int")
			}).Sql(),
			DownScript: NewDropTablesGenerator("test_table").Sql(),
		}
		Convey("Should apply migrations", func() {
			err = c.ApplyMigration(testMigration, false)
			So(err, ShouldBeNil)
			err = c.ApplyMigration(testMigration, true)
			So(err, ShouldBeNil)
		})
		Convey("Should insert migration", func() {
			So(c.InsertMigration(testMigration), ShouldBeNil)
			applied, err := c.GetAppliedMigrations()
			So(err, ShouldBeNil)
			So(applied, ShouldHaveLength, 1)
		})

		Convey("Should backup database", func() {
			_, err = c.Backup("/tmp/")
			So(err, ShouldBeNil)

		})
		Convey("Should remove migration", func() {
			So(c.RemoveMigration(testMigration), ShouldBeNil)
			applied, err := c.GetAppliedMigrations()
			So(err, ShouldBeNil)
			So(applied, ShouldHaveLength, 0)

		})
		Convey("Should remove migrations table", func() {
			c.DB.Exec(NewDropTablesGenerator("schema_version").Sql())
		})
	})
}
