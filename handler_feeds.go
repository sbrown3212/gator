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

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      feedName,
			Url:       feedURL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating new feed: %s", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following new feed: %s", err)
	}

	fmt.Println("Feed created successfully!")
	printFeed(feed)
	fmt.Printf("%s is now following \"%s\"", user.Name, feed.Name)

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

func handlerFollow(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current username: %s", err)
	}

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

	fmt.Printf("User %s is now following %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %s", err)
	}

	if len(cmd.Args) != 0 {
		return fmt.Errorf("%s command does not take any arguments", cmd.Name)
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting user's feeds: %s", err)
	}

	fmt.Printf("%s's feeds:\n", user.Name)
	for _, feed := range followedFeeds {
		if feed.FeedName.Valid {
			fmt.Printf(" * %s\n", feed.FeedName.String)
		} else {
			fmt.Println(" * (failed to get feed name)")
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
