package main

import (
	"fmt"
	"log"

	"github.com/sbrown3212/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("config: %+v\n", cfg)

	err = cfg.SetUser("stephen")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("config again: %+v\n", cfg)
}

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

// Run the given command with the current state, if it exists.
func (c *commands) run(s *state, cmd command) error {
	err := c.handlers[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

// Add a command to the commands struct.
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("please provide a username with the login command")
	}

	username := cmd.args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting username: %v", err)
	}

	fmt.Printf("User set to %s\n", username)
	return nil
}
