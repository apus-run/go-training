package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestRouter(t *testing.T) {
	// create web
	handler := Router()

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().ValueEqual("msg", "pong")

	obj := e.GET("/user/1").Expect().Status(http.StatusOK).JSON().Object()
	obj.Keys().ContainsOnly("user", "id")
	obj.ContainsKey("user").ValueEqual("user", "1")
	obj.Value("user").Object().ValueEqual("email", "foo@gmail.com")
}
