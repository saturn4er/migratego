package mysql

import (
"testing"

. "github.com/smartystreets/goconvey/convey"
)

func TestPrimaryKeyGenerator(t *testing.T) {
	Convey("Should create primary key without comment", t, func() {
		key := NewPrimaryKeyGenerator([]string{"a", "b"})
		So(key.comment, ShouldEqual, "")
		Convey("Should generate right sql", func() {
			So(key.Sql(), ShouldEqual, "PRIMARY KEY (`a`,`b`)")
		})
	})
	Convey("Should create primary key with comment ", t, func() {
		key := NewPrimaryKeyGenerator([]string{"a", "b", "`c"}, "test comment")
		So(key.comment, ShouldEqual, "test comment")
		Convey("Should generate right sql", func() {
			So(key.Sql(), ShouldEqual, "PRIMARY KEY (`a`,`b`,`\\`c`) COMMENT 'test comment'")
		})
	})
}
