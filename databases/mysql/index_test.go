package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIndexGenerator(t *testing.T) {
	Convey("Should create index generator", t, func() {
		idx := newIndexGenerator("idx_id", true).(*IndexGenerator)
		So(idx.name, ShouldEqual, "idx_id")
		So(idx.unique, ShouldEqual, true)
		idx.Name("idx_id1")
		So(idx.name, ShouldEqual, "idx_id1")
		idx.Comment("some comment")
		So(idx.comment, ShouldEqual, "some comment")
		idx.KeyBlockSize(100)
		So(idx.keyBlockSize, ShouldEqual, 100)
		idx.Parser("some parser")
		So(idx.parser, ShouldEqual, "some parser")
		So(idx.Sql(), ShouldBeEmpty)
		idx.Columns(&IndexColumnGenerator{
			Column: &ColumnGenerator{
				name: "id",
			},
			Order:  "ASC",
			Length: 10,
		})
		idx.Columns(&IndexColumnGenerator{
			Column: &ColumnGenerator{
				name: "id",
			},
			Order:  "DESC",
		})
		So(idx.columns, ShouldHaveLength, 2)
		So(idx.Sql(), ShouldEqual, "UNIQUE INDEX `idx_id1` (`id`(10) ASC,`id` DESC)")
	})
}
