package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sbrown3212/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <duration>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		scrapeFeeds(s.db)
	}
}

func scrapeFeeds(db *database.Queries) {
	// Query for next feed to fetch
	nextFeed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("error getting next feed to fetch:", err)
		return
	}

	fmt.Println("Found new feed to fetch...")
	scrapeFeed(db, nextFeed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now().UTC(),
		},
		ID: feed.ID,
	})
	if err != nil {
		log.Printf("unable to mark feed \"%s\" as fetched: %v\n", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("unable to fetch data for feed \"%s\": %s\n", feed.Name, err)
		return
	}

	fmt.Printf("Posts for feed \"%s\":\n", feed.Name)
	for _, item := range feedData.Channel.Item {
		fmt.Printf(" * %s\n", item.Title)
	}
	fmt.Printf("Collected %v posts for feed %s.\n", len(feedData.Channel.Item), feed.Name)

	fmt.Println()
}
