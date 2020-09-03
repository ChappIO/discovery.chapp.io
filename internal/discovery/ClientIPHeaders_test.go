package discovery

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestClientIPHeaders_GetClientIP(t *testing.T) {
	Convey("Given a ClientIPHeaders configured with 'X-One'", t, func() {
		headers := ClientIPHeaders{"X-One"}

		Convey("When no headers are present", func() {
			request := http.Request{
				RemoteAddr: "1.2.3.4:4433",
			}

			Convey("The return value is the remote ip address", func() {
				So(headers.GetClientIP(&request), ShouldEqual, "1.2.3.4")
			})
		})

		Convey("When the X-One header is set to a single ip", func() {
			request := http.Request{
				RemoteAddr: "1.2.3.4:4433",
				Header: http.Header{
					"X-One": {"1.1.1.1"},
				},
			}
			Convey("The return value is the header ip address", func() {
				So(headers.GetClientIP(&request), ShouldEqual, "1.1.1.1")
			})
		})

		Convey("When the X-One header is set to multiple ips", func() {
			request := http.Request{
				RemoteAddr: "1.2.3.4:4433",
				Header: http.Header{
					"X-One": {"2.1.1.1", "3.2.2.2"},
				},
			}
			Convey("The return value is the first value", func() {
				So(headers.GetClientIP(&request), ShouldEqual, "2.1.1.1")
			})
		})

		Convey("When the X-One header is set to an ip and port combo", func() {
			request := http.Request{
				RemoteAddr: "1.2.3.4:4433",
				Header: http.Header{
					"X-One": {"3.3.3.3:33"},
				},
			}
			Convey("The return value is the ip only", func() {
				So(headers.GetClientIP(&request), ShouldEqual, "3.3.3.3")
			})
		})
	})

}
