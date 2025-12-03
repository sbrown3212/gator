package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sbrown3212/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

// Run the given command with the current state, if it exists.
func (c *commands) run(s *state, cmd command) error {
	err := c.handlers[cmd.Name](s, cmd)
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
	if len(cmd.Args) == 0 {
		return fmt.Errorf("please provide a username with the login command")
	}

	username := cmd.Args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting username: %v", err)
	}

	fmt.Printf("User set to %s\n", username)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("config: %+v\n", cfg)

	programState := state{cfg: &cfg}

	handlers := make(map[string]func(*state, command) error)
	cmds := commands{handlers: handlers}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("please provide at least two arguments")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{Name: cmdName, Args: cmdArgs}

	err = cmds.run(&programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
