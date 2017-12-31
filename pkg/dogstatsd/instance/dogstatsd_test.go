package instance

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStart(t *testing.T) {
	Convey("Given an instance of dogstatsd", t, func() {
		statsd := NewDogStatsd()
		Convey("When the server is running", func() {
			statsd.Start()
			// Convey("Should not have any errors", func() {
			// So(err, ShouldBeNil)
			// })
			Convey("Should stop the server while a stop signal is given", func() {
				So(statsd.isStarted, ShouldBeTrue)
				statsd.Stop()
				So(statsd.isStarted, ShouldBeFalse)
			})
		})
	})
}
