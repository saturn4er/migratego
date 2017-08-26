package migratego

import (
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func FakeQB() QueryBuilder {
	return nil
}
func FakeDBClient(dsn, transactionTableName string) (DBClient, error) {
	return nil, errors.New("fake error")
}
func TestDrivers(t *testing.T) {
	Convey("Should add driver", t, func() {
		So(checkDriver("test_driver"), ShouldBeFalse)
		So(func() { getDriverClient("test_driver", "", "123") }, ShouldPanic)
		So(func() { getDriverQueryBuilder("test_driver") }, ShouldPanic)
		So(func() { shouldCheckDriver("test_driver") }, ShouldPanic)

		DefineDriver("test_driver", FakeQB, FakeDBClient)
		So(checkDriver("test_driver"), ShouldBeTrue)
		So(func() { shouldCheckDriver("test_driver") }, ShouldNotPanic)
		So(func() { getDriverClient("test_driver", "", "123") }, ShouldNotPanic)
		So(func() { getDriverQueryBuilder("test_driver") }, ShouldNotPanic)
	})
	Convey("Should not add second time driver with the same name", t, func() {
		So(func() { DefineDriver("test_driver", FakeQB, FakeDBClient) }, ShouldPanic)
	})
}
