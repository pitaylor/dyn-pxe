package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {

	t.Run("MimeType", func(t *testing.T) {
		var tpl Template

		tpl = Template{templateName: "template.txt"}
		assert.Equal(t, "text/plain; charset=utf-8", tpl.MimeType())

		tpl = Template{templateName: "template.html"}
		assert.Equal(t, "text/html; charset=utf-8", tpl.MimeType())

		tpl = Template{templateName: "template.unknown"}
		assert.Equal(t, "", tpl.MimeType())
	})

	t.Run("Render", func(t *testing.T) {
		tpl := newTemplate("test/template.txt")

		var err error
		var rendered string

		rendered, err = renderResource(tpl, ParamMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Template", rendered)

		rendered, err = renderResource(tpl, ParamMap{"p1": "v1"})
		assert.NoError(t, err)
		assert.Equal(t, "Template p1=v1", rendered)

		rendered, err = renderResource(tpl, ParamMap{"": "/template.txt"})
		assert.NoError(t, err)
		assert.Equal(t, "Template =/template.txt", rendered)
	})
}
