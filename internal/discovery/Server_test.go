package discovery

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func BodyShouldEqual(actual interface{}, expected ...interface{}) string {
	response := actual.(*httptest.ResponseRecorder)
	jsonBody := strings.TrimSpace(response.Body.String())
	expectedBody := strings.TrimSpace(expected[0].(string))
	if jsonBody != expectedBody {
		return fmt.Sprintf("Expected: [%s]\nFound:    [%s]", expectedBody, jsonBody)
	}
	return ""
}

func TestCoreServer_ServeHTTP(t *testing.T) {
	Convey("Given a preconfigured server'", t, func() {
		server := NewServer()

		Convey("When the root path is requested", func() {
			request := httptest.NewRequest("GET", "https://discovery.chapp.io", nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			Convey("I am redirected to the docs", func() {
				So(response.Code, ShouldEqual, http.StatusPermanentRedirect)
				So(response.Header().Get("location"), ShouldEqual, "https://github.com/ChappIO/discovery.chapp.io")
			})
		})

		Convey("When an agent is registered", func() {
			request := httptest.NewRequest("GET", "https://discovery.chapp.io/test.agent?private_address=1.1.2.2:223&agent_id=23rf42ef2", nil)
			request.Header.Set("x-forwarded-for", "1.1.1.1")
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			Convey("Its details are returned", func() {
				So(response, BodyShouldEqual, `{"service_id":"test.agent","public_ip":"1.1.1.1","agents":[{"agent_id":"23rf42ef2","private_address":"1.1.2.2:223"}]}`)
			})
		})

		Convey("When a non-existing agent is requested", func() {
			request := httptest.NewRequest("GET", "https://discovery.chapp.io/test.empty", nil)
			request.Header.Set("x-forwarded-for", "1.1.1.1")
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			Convey("An empty list is returned", func() {
				So(response, BodyShouldEqual, `{"service_id":"test.empty","public_ip":"1.1.1.1","agents":[]}`)
			})
		})

		Convey("When an existing agent is requested", func() {
			// preload
			server.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "https://discovery.chapp.io/test.exists?private_address=1.2.3.4:567&agent_id=jupgottem", nil))
			server.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "https://discovery.chapp.io/test.exists?private_address=1.2.3.4:567&agent_id=jupgottem2", nil))

			request := httptest.NewRequest("GET", "https://discovery.chapp.io/test.exists", nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			Convey("Its details are returned", func() {
				So(response, BodyShouldEqual, `{"service_id":"test.exists","public_ip":"192.0.2.1","agents":[{"agent_id":"jupgottem","private_address":"1.2.3.4:567"},{"agent_id":"jupgottem2","private_address":"1.2.3.4:567"}]}`)
			})
		})

		Convey("When an agent id is too long", func() {
			request := httptest.NewRequest("GET", "https://discovery.chapp.io/long.agent?private_address=1.1.2.2:223&agent_id=1234567890123456789012345678901234567890", nil)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			Convey("Its id is shortened to 36 characters", func() {
				So(response, BodyShouldEqual, `{"service_id":"long.agent","public_ip":"192.0.2.1","agents":[{"agent_id":"123456789012345678901234567890123456","private_address":"1.1.2.2:223"}]}`)
			})
		})
	})
}
