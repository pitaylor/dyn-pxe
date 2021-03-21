package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Command is a resource that represents an external program.
type Command struct {
	command string
}

func newCommand(cmd string) *Command {
	return &Command{command: cmd}
}

// MimeType is implemented to satisfy the Resource interface but returns an
// empty string for now.
func (c *Command) MimeType() string {
	return ""
}

// Render executes the command and outputs the command's stdout. Params are
// passed to the command as environment variables.
func (c *Command) Render(out io.Writer, params ParamMap, vars VariableMap) error {
	env := make([]string, len(params))
	i := 0

	for k, v := range params {
		if k != "" {
			env[i] = fmt.Sprintf("%v=%v", k, v)
		}
		i++
	}

	for k, v := range vars {
		if k != "" {
			env[i] = fmt.Sprintf("%v=%v", k, v)
		}
	}

	split := strings.Split(c.command, " ")
	cmd := exec.Command(split[0], split[1:]...)
	cmd.Env = append(os.Environ(), env...)

	log.Printf("Render command: %v %v", c.command, env)
	output, err := cmd.Output()
	out.Write(output)

	return err
}
