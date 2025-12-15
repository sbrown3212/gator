package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sbrown3212/gator/internal/database"
)

const (
	longForm  = "Jan 2, 2006 at 3:04pm (MST)"
	shortForm = "2006-Jan-02"
	RFC3339   = time.RFC3339
	RFC1123Z  = time.RFC1123Z
	RFC1123   = time.RFC1123
)

var timeFormats = []string{
	RFC1123Z,
	RFC1123,
	RFC3339,
	longForm,
	shortForm,
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <duration>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n\n", timeBetweenReqs)

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

	log.Println("Found new feed to fetch...")
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
		log.Printf("unable to fetch data for feed \"%s\": %s\n\n", feed.Name, err)
		return
	}

	log.Printf("Feed: %s:\n", feed.Name)

	count := 0

	for _, item := range feedData.Channel.Item {
		descriptionParameter := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		pubDateParameter, dateErr := getSQLNullTime(item.PubDate)

		feed, err := db.GetFeedByUrl(context.Background(), feed.Url)
		if err != nil {
			log.Printf("unable to find feed for post %s", item.Title)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: descriptionParameter,
			PublishedAt: pubDateParameter,
			FeedID:      feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					// Ignore (post already exists in db)
					continue
				}
			}

			log.Printf(" - error (Feed: %s, Post: %s): %v", feed.Name, item.Title, err)
			continue
		}

		log.Printf(" * (ok) %s", item.Title)

		if dateErr != nil {
			log.Printf("   - %s", dateErr)
		}

		count++
	}
	log.Printf("Found %v new posts for feed %s", count, feed.Name)
	fmt.Println()
}

func getSQLNullTime(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{
			Valid: false,
		}, fmt.Errorf("date string is empty")
	}

	for _, format := range timeFormats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return sql.NullTime{
				Valid: true,
				Time:  t,
			}, nil
		}
	}

	return sql.NullTime{
		Valid: false,
	}, fmt.Errorf("unable to parse date string")
}
