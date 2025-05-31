package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iegpeppino/blog_aggregator/internal/database"
)

// Handles RSSFeed get request
func handleAgg(s *state, cmd command) error {

	// Queries all rssFeed items from the provided url
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch rssfeed %w", err)
	}

	//fmt.Printf("Feed %+v\n", rssFeed)
	for _, item := range rssFeed.Channel.Item[:1] {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
	}

	return nil
}

// Handles adding a new feed into db
func handleAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 2 {
		return errors.New("invalid arguments") // Needs two arguments or raises error
	}

	ctx := context.Background()
	feedName := cmd.Args[0] // Store argument values
	feedURL := cmd.Args[1]

	// Set the new rssFeed parameters
	// using the UUID from the current user (foreign key)
	feedParams := database.CreateFeedParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	// Queries the insertion of the feed item into the db
	rssFeed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("couldn't create feed %w", err)
	}
	fmt.Println("Succesfully created new feed:")
	fmt.Printf("Title: %v\n", rssFeed.Name)
	fmt.Printf("URL: %v\n", rssFeed.Url)

	// Creates a feed follow between the user and it's
	// new feed
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    rssFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed follow created succesfully!")

	return nil
}

// Retrieves all feeds and it's user data
func handleFeeds(s *state, cmd command) error {

	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get feeds %w", err)
	}

	// Get the user data for every feed in feeds
	for _, feed := range feeds {
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user data %w", err)
		}

		fmt.Println(feed.Name) // Print results
		fmt.Println(feed.Url)
		fmt.Println(user.Name)
	}

	return nil
}

// Deletes a feed_follow record  (aka. unfollows a feed)
func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("invalid arguments, only feed URL is required")
	}

	// Get context and feedURL
	ctx := context.Background()
	feedURL := cmd.Args[0]

	// Get feed to access it's ID
	feed, err := s.db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return err
	}

	// Set parameters for DELETE query
	delParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	// Performs query
	err = s.db.DeleteFeedFollow(ctx, delParams)
	if err != nil {
		return fmt.Errorf("couldn't unfollow %w", err)
	}

	return nil
}
