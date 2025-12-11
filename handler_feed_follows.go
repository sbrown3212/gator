package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sbrown3212/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("unable to find feed (create feed before following): %s", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %s", err)
	}

	fmt.Printf("Feed followed successfully!\n")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("%s command does not take any arguments", cmd.Name)
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting user's feeds: %s", err)
	}

	fmt.Printf("%s's feeds:\n", user.Name)
	for _, feed := range followedFeeds {
		fmt.Printf(" * %s\n", feed.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf(" * User: %s\n", username)
	fmt.Printf(" * Feed: %s\n", feedname)
}
