package public

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jcarley/datica-users/web"
	"github.com/jcarley/s3lite/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHealthCheckGreeting(t *testing.T) {

	controller := NewHealthCheckController()

	Convey("Outputs hello world greeting with unauthenticated request", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		rw := httptest.NewRecorder()

		controller.Show(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		greeting := rawData["greeting"]

		So(greeting, ShouldEqual, "Hello World")
	})
}

func TestHealthCheckGreeting_Integration(t *testing.T) {

	server := web.NewHttpServer()

	controller := NewHealthCheckController()
	controller.Register(server)

	Convey("Outputs hello world greeting with unauthenticated request", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		rw := httptest.NewRecorder()

		server.Router.ServeHTTP(rw, req)

		buffer := rw.Body.Bytes()
		rawData := test.GetRawData(t, buffer)

		greeting := rawData["greeting"]

		So(greeting, ShouldEqual, "Hello World")
	})
}
