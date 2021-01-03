package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd := newCommand("test/command.sh")

	t.Run("MimeType", func(t *testing.T) {
		assert.Equal(t, "", cmd.MimeType())
	})

	t.Run("Render", func(t *testing.T) {
		var err error
		var rendered string

		rendered, err = renderResource(cmd, ParamMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh", rendered)

		rendered, err = renderResource(cmd, ParamMap{"p1": "v1", "p2": "v2"})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh p1=v1 p2=v2", rendered)

		rendered, err = renderResource(cmd, ParamMap{"": "/command.sh"})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh", rendered)

		rendered, err = renderResource(newCommand("test/command.sh 1 2 3"), ParamMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh \"1\" \"2\" \"3\"", rendered)

		rendered, err = renderResource(cmd, ParamMap{"EXITCODE": "1"})
		assert.Error(t, err)
		assert.Equal(t, "Command: test/command.sh", rendered)
	})
}
