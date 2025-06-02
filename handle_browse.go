package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/iegpeppino/blog_aggregator/internal/database"
)

// Gets posts from feeds the current user follows and prints
// their contents on the terminal
func handleBrowse(s *state, cmd command, user database.User) error {

	ctx := context.Background()

	// If no limit arg is passed, it defaults to 2
	limit := 2
	if len(cmd.Args) == 1 {
		if limitInput, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = limitInput
		} else {
			return fmt.Errorf("invalid argument: %w", err)
		}
	}

	// Set query parameters
	getPostsParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	// Query posts from table
	posts, err := s.db.GetPostsForUser(ctx, getPostsParams)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	// Printing posts' contents
	fmt.Printf("Got %v posts for user %s :\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Println("Printing post")
		fmt.Printf("- Feed: %s (%s)\n", post.FeedName, post.PublishedAt.Time.Format("Mon Jan 02"))
		fmt.Printf("-- %s --\n", post.Title)
		fmt.Printf("	%v\n", post.Description.String)
		fmt.Printf("Link to post: %s\n", post.Url)
		fmt.Println("######################################################")
	}

	return nil
}
