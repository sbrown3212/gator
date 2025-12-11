package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting users: %w", err)
	}

	fmt.Println("Successfully reset users.")
	return nil
}
