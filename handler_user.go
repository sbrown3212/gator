package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting username: %v", err)
	}

	fmt.Printf("User set to %s\n", username)
	return nil
}
