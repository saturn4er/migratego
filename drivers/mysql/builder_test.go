package mysql

import (
	"testing"

	"github.com/saturn4er/migratego"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMysqlQueryBuilder_CreateTable(t *testing.T) {

	Convey("QueryBuilder.CreateTable Should generate right sql", t, func() {
		b := MysqlQueryBuilder{}
		b.CreateTable("test_table", func(g migratego.CreateTableGenerator) {
			id := g.Column("id", "varchar(255)")
			So(id.GetName(), ShouldEqual, "id")
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) NULL)")
			id.Binary()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NULL)")
			id.Comment("test_comment")
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NULL COMMENT 'test_comment')")
			id.DefaultValue("default_value")
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NULL DEFAULT 'default_value' COMMENT 'test_comment')")
			id.Generated()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NULL GENERATED ALWAYS AS ('default_value') COMMENT 'test_comment')")
			id.NotNull()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NOT NULL GENERATED ALWAYS AS ('default_value') COMMENT 'test_comment')")
			id.Primary()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NOT NULL GENERATED ALWAYS AS ('default_value') COMMENT 'test_comment', PRIMARY KEY (`id`))")
			id.Primary("primary_comment")
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) BINARY NOT NULL GENERATED ALWAYS AS ('default_value') COMMENT 'test_comment', PRIMARY KEY (`id`) COMMENT 'primary_comment')")
			id.ZeroFill()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) ZEROFILL BINARY NOT NULL GENERATED ALWAYS AS ('default_value') COMMENT 'test_comment', PRIMARY KEY (`id`) COMMENT 'primary_comment')")
			id.Unsigned()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) UNSIGNED ZEROFILL BINARY NOT NULL GENERATED ALWAYS AS ('default_value') COMMENT 'test_comment', PRIMARY KEY (`id`) COMMENT 'primary_comment')")
			id.AutoIncrement()
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) UNSIGNED ZEROFILL BINARY NOT NULL GENERATED ALWAYS AS ('default_value') AUTO_INCREMENT COMMENT 'test_comment', PRIMARY KEY (`id`))")
			id.AutoIncrement("autoincrement_comment")
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) UNSIGNED ZEROFILL BINARY NOT NULL GENERATED ALWAYS AS ('default_value') AUTO_INCREMENT COMMENT 'test_comment', PRIMARY KEY (`id`) COMMENT 'autoincrement_comment')")
			id.Index("idx_id", true)
			So(g.Sql(), ShouldEqual, "CREATE TABLE `test_table`(`id` varchar(255) UNSIGNED ZEROFILL BINARY NOT NULL GENERATED ALWAYS AS ('default_value') AUTO_INCREMENT COMMENT 'test_comment', PRIMARY KEY (`id`) COMMENT 'autoincrement_comment', UNIQUE INDEX `idx_id` (`id` ASC))")
		})
		g := b.DropTables("test_table").IfExists()
		So(g.Sql(), ShouldEqual, "DROP TABLE IF EXISTS `test_table`")
		g.Table("test_table_2")
		So(g.Sql(), ShouldEqual, "DROP TABLE IF EXISTS `test_table`,`test_table_2`")
		g = b.DropTables().IfExists()
		So(g.Sql(), ShouldEqual, "")
	})
}
