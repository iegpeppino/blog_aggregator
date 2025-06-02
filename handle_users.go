package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iegpeppino/blog_aggregator/internal/database"
)

// Logs in a user (if it exists in the users table)
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("invalid argument/s\n Usage: %v <username>", cmd.Name)
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

// Inserts a new user record into users table
func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("invalid argument/s\nUsage: %v <username>", cmd.Name)
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
	fmt.Printf(" -User Name: %v\n", userParams.Name)

	return nil
}

// Gets all users from users table and prints their names
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

// Deletes all records from users table
func handleReset(s *state, cmd command) error {

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset database %w", err)
	}

	fmt.Println("Database succesfully reset!")

	return nil
}
