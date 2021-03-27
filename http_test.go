package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

		data := url.Values{}
		data.Set("v1", "post")

		req, err := http.NewRequest("POST", "/command/x/y?v2=query", strings.NewReader(data.Encode()))
		assert.NoError(t, err)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
		assert.Equal(t, "Command: test/command.sh p1=x p2=y v1=post v2=query HTTP_METHOD=POST", rr.Body.String())
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
