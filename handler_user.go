package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sbrown3212/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("user not found: %s", err)
	}

	err = s.cfg.SetUser(user.Name)
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

	ctx := context.Background()

	queryArgs := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: username,
	}
	user, err := s.db.CreateUser(ctx, queryArgs)
	if err != nil {
		return fmt.Errorf("error adding user to database: %s", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %s", err)
	}

	fmt.Println("New user registered!")
	log.Printf("New user: %v", user)

	return nil
}
