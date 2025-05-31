package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iegpeppino/blog_aggregator/internal/database"
)

// Creates a new feed_follow record
func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("invalid arguments, please provide only an url")
	}

	ctx := context.Background()

	feedURL := cmd.Args[0]

	// Query the desired feed in order to get its id
	feed, err := s.db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return err
	}

	// Set feed_follow parameters
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	// Create new feed_follow record
	_, err = s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return fmt.Errorf("couldn't create feed_follow %w", err)
	}

	fmt.Printf("User %s now follows %s ", user.Name, feedURL)

	return nil
}

// Returns all followed RRSFeeds from a user
func handleFollowing(s *state, cmd command, user database.User) error {

	ctx := context.Background()

	// Get feeds followed by our user
	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get followed feeds %w", err)
	}

	// Print results
	fmt.Printf("User %s follows these feeds:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf(" - %s\n", follow.FeedName)
	}

	return nil
}
