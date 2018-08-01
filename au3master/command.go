package au3master

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Command is a representation of a autoit 3 function call
type Command struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
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
func (cmd *Command) SetParams(v ...interface{}) {
	p := make([]interface{}, len(v))
	for k, val := range v {
		p[k] = val
	}
	cmd.Params = p
}

// Result is the result of an autoit command
type Result struct {
	CommandID string
	Value     json.RawMessage
}

type au3resp struct {
	Success   bool            `json:"success"`
	Type      string          `json:"type"`
	CommandID string          `json:"command_id,omitempty"`
	Value     json.RawMessage `json:"value,omitempty"`
}
