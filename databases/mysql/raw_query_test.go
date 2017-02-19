package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRawQuery(t *testing.T) {
	Convey("Should return same query for RawQuery", t, func() {
		q := "abcdefghijklmnopqrstu"
		rawQ := rawQuery(q)
		So((&rawQ).Sql(), ShouldEqual, q)
	})
}
