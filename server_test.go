package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestServer(t *testing.T) {
	t.Run("UnmarshalYAML", func(t *testing.T) {
		y := normalizeYaml(`
			http-address: 1.1.1.1:8080
			tftp-address: 2.2.2.2:169
			static-dir: /var/static

			resources:
			  - route: /path/:p1/command.sh
				command: test/command.sh

			  - route: /path/(?P<p1>\w+)/template.txt
			    template: test/template.txt
		`)

		s := &Server{}
		err := yaml.Unmarshal(y, &s)
		assert.NoError(t, err)
		assert.Equal(t, "/var/static", s.StaticDir)
		assert.Equal(t, "1.1.1.1:8080", s.HTTPAddress)
		assert.Equal(t, "2.2.2.2:169", s.TFTPAddress)

		rendered, err := renderRoute(s.Routes, "/path/123/command.sh")
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh p1=123", rendered)

		rendered, err = renderRoute(s.Routes, "/path/456/template.txt")
		assert.NoError(t, err)
		assert.Equal(t, "Template =/path/456/template.txt p1=456", rendered)

		y = normalizeYaml(`
			resources:
			  - command: test/command.sh
		`)

		err = yaml.Unmarshal(y, &s)
		assert.Error(t, err)
		assert.Equal(t, "resource[0]: missing route", err.Error())

		y = normalizeYaml(`
			resources:
			  - route: /xyz
		`)

		err = yaml.Unmarshal(y, &s)
		assert.Error(t, err)
		assert.Equal(t, "resource[0]: exactly one of `command` or `template` expected", err.Error())
	})
}
