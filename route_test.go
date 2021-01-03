package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRoute(t *testing.T) {
	t.Run("UnmarshalYAML", func(t *testing.T) {
		const y = `
			routes:	
			- /a
			- /b/:p1
		`

		r := struct{ Routes RouteSet }{}
		err := yaml.Unmarshal(normalizeYaml(y), &r)
		assert.NoError(t, err)
		assert.Equal(t, "/a", r.Routes[0].Path)
		assert.Equal(t, "/b/:p1", r.Routes[1].Path)
	})
}

func TestRouteSet(t *testing.T) {
	t.Run("match", func(t *testing.T) {
		seq := RouteSet{
			newRoute("/a", nil),
			newRoute("/b/:p1", nil),
		}

		var r *Route
		var p ParamMap

		r, p = seq.match("x")
		assert.Nil(t, r)

		r, p = seq.match("/a")
		assert.Equal(t, seq[0], r)
		assert.Equal(t, ParamMap{"": "/a"}, p)

		r, p = seq.match("/b/c")
		assert.Equal(t, seq[1], r)
		assert.Equal(t, ParamMap{"": "/b/c", "p1": "c"}, p)
	})
}
