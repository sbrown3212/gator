package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sbrown3212/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed db entry: %s", err)
	}

	fmt.Println("Feed created successfully!")
	printFeed(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("%s takes 0 arguments", cmd.Name)
	}

	feeds, err := s.db.GetFeedsAndUsername(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Found %v feeds:\n", len(feeds))
	for i, feed := range feeds {
		if i != 0 {
			fmt.Println("---")
		}
		fmt.Printf(" * Name:       %s\n", feed.Name)
		fmt.Printf(" * URL:        %s\n", feed.Url)
		if feed.User.Valid {
			fmt.Printf(" * Created by: %s\n", feed.User.String)
		} else {
			fmt.Println(" * Created by: not listed")
		}
	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %s\n", feed.ID)
	fmt.Printf(" * Created: %v\n", feed.CreatedAt)
	fmt.Printf(" * Updated: %v\n", feed.UpdatedAt)
	fmt.Printf(" * Name:    %s\n", feed.Name)
	fmt.Printf(" * URL:     %s\n", feed.Url)
	fmt.Printf(" * UserID:  %s\n", feed.UserID)
}
