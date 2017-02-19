package migratego

import (
	"testing"

	"sort"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMigration(t *testing.T) {
	Convey("Should compare migrations", t, func() {
		a := &Migration{
			Number:     1,
			Name:       "test",
			UpScript:   "SELECT 1",
			DownScript: "SELECT 2",
		}
		b := &Migration{
			Number:     1,
			Name:       "test",
			UpScript:   "SELECT 1",
			DownScript: "SELECT 2",
		}
		So(a.Compare(b), ShouldBeTrue)
		now := time.Now()
		b.AppliedAt = &now
		So(a.Compare(b), ShouldBeTrue)
		b.Number = 10
		So(a.Compare(b), ShouldBeFalse)
	})
	Convey("Should test Migration.Compare", t, func() {
		a := []Migration{
			{
				Number:     3,
				Name:       "3",
			},
			{
				Number:     1,
				Name:       "1",
			},
			{
				Number:     2,
				Name:       "2",
			},
		}
		sort.Sort(byNumber(a))
		So(a[0].Name, ShouldEqual, "1")
		So(a[1].Name, ShouldEqual, "2")
		So(a[2].Name, ShouldEqual, "3")
	})
}
