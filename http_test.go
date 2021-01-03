package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPHandler(t *testing.T) {
	s := Server{
		StaticDir: "test",
		Routes: RouteSet{
			newRoute("/command/:p1/:p2", newCommand("test/command.sh")),
			newRoute("/template/:p1/:p2", newTemplate("test/template.txt")),
		},
	}

	handler := s.newHTTPHandler()

	t.Run("template", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/template/x/y", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
		assert.Equal(t, "Template =/template/x/y HTTP_METHOD=GET p1=x p2=y", rr.Body.String())
	})

	t.Run("does-not-exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/not-found", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "404 page not found\n", rr.Body.String())
	})

	t.Run("command", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/command/x/y", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
		assert.Equal(t, "Command: test/command.sh p1=x p2=y HTTP_METHOD=POST", rr.Body.String())
	})

	t.Run("static-file", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/hello.txt", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
		assert.Equal(t, "Hello World\n", rr.Body.String())
	})
}
