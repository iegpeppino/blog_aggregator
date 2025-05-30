package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iegpeppino/blog_aggregator/internal/database"
)

// Logs in a user (if it exists in the users DB)
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("login argument missing or multiple arguments present")
	}

	userName := cmd.Args[0]

	// Checking if user exists
	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("couldn't find user %w", err)
	}
	// Setting user to the entered username in args
	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("couldn't set user %w", err)
	}

	fmt.Printf("Username '%v' has been set", cmd.Args[0])

	return nil
}

// Inserts a new user record into users DB
func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("register command missing an argument")
	}

	userName := cmd.Args[0]

	// Setting user params to be passed on
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userName,
	}

	// Calling the query
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("couldn't create user entry %w", err)
	}

	// Setting recently created user as current user
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user %w", err)
	}

	fmt.Println("User created succesfully!")
	fmt.Printf(" -Id: %v\n", userParams.ID)
	fmt.Printf(" -User Name: %v\n", userParams.Name)

	return nil
}

// Erases all records from users DB
func handleReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset database %w", err)
	}
	fmt.Println("Database succesfully reset!")
	return nil
}

// Gets all users from DB and prints their names
func handleGetUsers(s *state, cmd command) error {

	// Queries all users
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list all users %w", err)
	}

	// Iterates the slice of users and prints their names
	// tags the loged user as (current)
	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}
	return nil
}

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

	// Needs two arguments or raises error
	if len(cmd.Args) != 2 {
		return errors.New("invalid arguments")
	}

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
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	// Queries the insertion of the feed item into the db
	rssFeed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("couldn't create feed %w", err)
	}

	fmt.Printf("Title: %v\n", rssFeed.Name)
	fmt.Printf("URL: %v\n", rssFeed.Url)

	return nil
}

func handleFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get feeds %w", err)
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user data %w", err)
		}
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(user.Name)
	}
	return nil
}
