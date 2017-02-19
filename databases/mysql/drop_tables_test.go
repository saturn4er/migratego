package mysql

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDropTablesGenerator(t *testing.T) {
	Convey("Should create drop tables generator", t, func() {
		tableNames := []string{"a", "b", "c"}
		g := NewDropTablesGenerator(tableNames...).(*dropTablesGenerator)
		So(g.ifExists, ShouldBeFalse)
		So(g.tables, should.HaveLength, len(tableNames))
		g.Table("d")
		So(g.tables, should.HaveLength, len(tableNames)+1)
		So(g.Sql(), ShouldEqual, "DROP TABLE `a`,`b`,`c`,`d`")
		g.IfExists()
		So(g.Sql(), ShouldEqual, "DROP TABLE IF EXISTS `a`,`b`,`c`,`d`")
		g.tables = []string{}
		So(g.Sql(), ShouldEqual, "")
	})
}
