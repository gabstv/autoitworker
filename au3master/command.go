package au3master

import (
	"github.com/google/uuid"
)

// Command is a representation of a autoit 3 function call
type Command struct {
	ID     string
	Name   string
	Params []string
}

func newCommand(name string) *Command {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	cmd := &Command{
		ID:   uid.String(),
		Name: name,
	}
	return cmd
}

// SetParams assigns the Params of the command struct
func (cmd *Command) SetParams(v ...string) {
	p := make([]string, len(v))
	for k, val := range v {
		p[k] = val
	}
	cmd.Params = p
}

// Result is the result of an autoit command
type Result struct {
	CommandID string
}
