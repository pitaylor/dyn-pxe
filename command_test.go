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

		rendered, err = renderResource(cmd, ParamMap{}, VariableMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh", rendered)

		rendered, err = renderResource(cmd, ParamMap{"p1": "v1", "p2": "v2"}, VariableMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh p1=v1 p2=v2", rendered)

		rendered, err = renderResource(cmd, ParamMap{"p1": "x"}, VariableMap{"v1": "y"})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh p1=x v1=y", rendered)

		rendered, err = renderResource(cmd, ParamMap{"": "/command.sh"}, VariableMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh", rendered)

		rendered, err = renderResource(newCommand("test/command.sh 1 2 3"), ParamMap{}, VariableMap{})
		assert.NoError(t, err)
		assert.Equal(t, "Command: test/command.sh \"1\" \"2\" \"3\"", rendered, VariableMap{})

		rendered, err = renderResource(cmd, ParamMap{"EXITCODE": "1"}, VariableMap{})
		assert.Error(t, err)
		assert.Equal(t, "Command: test/command.sh", rendered)
	})
}
