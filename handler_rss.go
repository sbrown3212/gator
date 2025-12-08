package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	// if len(cmd.Args) != 1 {
	// 	return fmt.Errorf("usage: %s <rss-feed-url>", cmd.Name)
	// }

	// feedURL := cmd.Args[0]

	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %s", err)
	}

	// printFeed(feed)
	fmt.Print(feed)
	return nil
}

// func printFeed(feed *RSSFeed) {
// 	fmt.Printf("Channel Title: %s\n", feed.Channel.Title)
// 	fmt.Printf("Channel Description: %s\n\n", feed.Channel.Description)
//
// 	for i, item := range feed.Channel.Item {
// 		fmt.Printf("%v:\n", i+1)
// 		fmt.Printf(" - Title: %s\n", item.Title)
// 		fmt.Printf(" - Description: %s\n", item.Description)
// 	}
// }
