package main

import (
	"context"
	"errors"
	"fmt"
	"time"

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
func handleAddFeed(s *state, cmd command) error {

	if len(cmd.Args) != 2 {
		return errors.New("invalid arguments") // Needs two arguments or raises error
	}

	feedName := cmd.Args[0] // Store argument values
	feedURL := cmd.Args[1]

	// Gets the user data from the current user
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

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
	rssFeed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("couldn't create feed %w", err)
	}

	fmt.Printf("Title: %v\n", rssFeed.Name)
	fmt.Printf("URL: %v\n", rssFeed.Url)

	// Calls the follow handler for the feed we've just created

	cmd.Args = []string{feedURL} // Modify arguments to be able to call handleFollow

	err = handleFollow(s, cmd)
	if err != nil {
		return fmt.Errorf("couldn't follow created feed %w", err)
	}

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
