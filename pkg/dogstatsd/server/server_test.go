package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRun(t *testing.T) {
	Convey("Given an instance of server", t, func() {
		srv := NewServer()
		Convey("When server receives metrics in JSON format", func() {
			w := httptest.NewRecorder()
			data := []byte(`{"series": [{
				"source_type_name": "System",
				"metric": "ntp.offset",
				"points": [
					[
						1514917995,
						0.002307891845703125
					]
				],
				"type": "gauge",
				"host": "hsiang-ubuntu"
				}
			]}`)
			req, _ := http.NewRequest("POST", "/v1/dogstatsd", bytes.NewBuffer(data))
			srv.Handler.ServeHTTP(w, req)
			Convey("Should be able to process the request", func() {
				So(string(w.Body.Bytes()[:]), ShouldEqual, `{"response":"received"}`)
				So(w.Code, ShouldEqual, 200)
			})
		})

		Convey("When server is running", func() {
			Convey("Should not throw an error if the server is being initialized", func() {
				var err error
				var wg sync.WaitGroup
				wg.Add(1)
				ch := make(chan error)

				go func() {
					wg.Done()
					srv.Run(ch)
				}()
				wg.Wait()

				So(srv.isInitialized, ShouldBeTrue)
				So(err, ShouldBeNil)

				go srv.Run(ch)
				select {
				case msg := <-ch:
					err = msg
				case <-time.After(time.Second * 1):
					err = nil
				}
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, HasBeenInitialized)

				So(srv.Stop(), ShouldBeNil)
				So(srv.isInitialized, ShouldBeFalse)
			})
		})
	})
}
