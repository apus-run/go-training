package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"

	"gin-with-seesion/internal/web/handler"
)

func TestRouter(t *testing.T) {
	userHandler := handler.NewUserHandler()
	// create web
	router := Router(userHandler)

	// run server using httptest
	server := httptest.NewServer(router)
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
