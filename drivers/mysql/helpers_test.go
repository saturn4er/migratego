package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelpers(t *testing.T) {
	Convey("Should wrap name", t, func() {
		So(wrapName("`hello`qwe"), ShouldEqual, "`\\`hello\\`qwe`")
		So(wrapNames([]string{"a", "`hello`qwe", "b"}), ShouldEqual, "`a`,`\\`hello\\`qwe`,`b`")
	})
}
