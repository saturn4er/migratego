package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIndexColumnGenerator(t *testing.T) {
	Convey("Should create index column generator", t, func() {
		idx := NewIndexColumnGenerator("hello", "ASC", 10).(*IndexColumnGenerator)
		So(idx.Length, ShouldEqual, 10)
		So(idx.Order, ShouldEqual, "ASC")
		So(idx.Sql(), ShouldEqual, "`hello`(10) ASC")
	})
}
