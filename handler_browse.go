package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sbrown3212/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	}

	limit := int32(2)
	if len(cmd.Args) == 1 {
		argAsNum, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("please ensure %s command argument can be parsed as an integer", cmd.Name)
		}
		limit = int32(argAsNum)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}

	fmt.Printf("Here are your %v most recent posts:\n", limit)
	fmt.Println("(navigate to a post by holding the command key and clicking that post's link)")

	for _, post := range posts {
		printPost(post)
		fmt.Println()
	}

	return nil
}

func printPost(post database.GetPostsForUserRow) {
	pubDate := "(not listed)"
	if post.PublishedAt.Valid {
		pubDate = post.PublishedAt.Time.Format("Mon Jan 2")
	}

	fmt.Printf("%s from %s\n", pubDate, post.FeedName)
	fmt.Printf("--- %s ---\n", post.Title)
	fmt.Printf("    %s\n", post.Description.String)
	fmt.Printf("Link: %s\n", post.Url)
}
