package message

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPush(t *testing.T) {
	Convey("When Push function is called", t, func() {
		expect := "TEST_DATA"
		Push([]byte(expect))
		Convey("Should be able to retrive data from Data()", func() {
			metrics := <-Data()
			So(string(metrics[:]), ShouldEqual, expect)
		})
	})
}
