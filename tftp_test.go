package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testReader struct {
	output string
}

func (t *testReader) ReadFrom(r io.Reader) (int64, error) {
	buf := make([]byte, 1024)
	bytesRead := int64(0)

	for {
		n, err := r.Read(buf)
		bytesRead += int64(n)
		t.output += string(buf[:n])

		if err == io.EOF {
			break
		} else if err != nil {
			return bytesRead, err
		}
	}

	return bytesRead, nil
}

func TestNewTFTPReadHandler(t *testing.T) {
	s := Server{
		StaticDir: "test",
		Routes: RouteSet{
			newRoute("/command/:p1/:p2", newCommand("test/command.sh")),
			newRoute("/template/:p1/:p2", newTemplate("test/template.txt")),
		},
	}

	handler := s.newTFTPReadHandler()

	t.Run("template", func(t *testing.T) {
		r := &testReader{}
		err := handler("/template/x/y", r)
		assert.NoError(t, err)
		assert.Equal(t, "Template =/template/x/y p1=x p2=y", r.output)
	})

	t.Run("does-not-exist", func(t *testing.T) {
		r := &testReader{}
		err := handler("/not-found", r)
		assert.Error(t, err)
		assert.Equal(t, "", r.output)
	})

	t.Run("static-file", func(t *testing.T) {
		r := &testReader{}
		err := handler("/hello.txt", r)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World\n", r.output)
	})

	t.Run("static-file-unsafe", func(t *testing.T) {
		r := &testReader{}
		err := handler("../../hello.txt", r)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World\n", r.output)
	})
}
