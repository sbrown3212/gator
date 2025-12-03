package main

import (
	"log"
	"os"

	"github.com/sbrown3212/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{cfg: &cfg}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("please provide at least two arguments")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
