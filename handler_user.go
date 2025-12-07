package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sbrown3212/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not found: %s", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting username: %v", err)
	}

	fmt.Printf("User set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	queryArgs := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		Name: username,
	}

	user, err := s.db.CreateUser(context.Background(), queryArgs)
	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %s", err)
	}

	fmt.Println("New user registered!")
	printUser(user)

	return nil
}

func handlerUser(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error looking up users: %s", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("there are currently no registered users")
	}

	for _, user := range users {
		display := user.Name

		if user.Name == s.cfg.CurrentUserName {
			display += " (current)"
		}

		fmt.Printf(" * %s\n", display)
	}
	return nil
}

func printUser(user database.User) {
	fmt.Println(" * ID:", user.ID)
	fmt.Println(" * Name:", user.Name)
}
