package cli

import "fmt"

type Commands struct {
	handlers map[string]HandlerFunc
}

func NewCommands() *Commands {
	return &Commands{handlers: map[string]HandlerFunc{}}
}

func (c *Commands) Register(name string, f HandlerFunc) {
	c.handlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	h, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return h(s, cmd)
}
