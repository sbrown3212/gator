package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// Run the given command with the current state, if it exists.
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	return f(s, cmd)
}

// Add a command to the commands struct.
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
