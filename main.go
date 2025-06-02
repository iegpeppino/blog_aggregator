package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/iegpeppino/blog_aggregator/internal/config"
	"github.com/iegpeppino/blog_aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {

	// Reads the configuration json file to set it
	// into the program state
	configStruct, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	// Load the DB URL and open a connection to the database
	db, err := sql.Open("postgres", configStruct.DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Create new db and predetermined queries
	// made by sqlc
	dbQueries := database.New(db)

	programState := &state{
		cfg: &configStruct,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCmds: make(map[string]func(*state, command) error),
	}

	// Registering our custom CLI commands

	cmds.register("login", handlerLogin)                            // Logs in user
	cmds.register("register", handleRegister)                       // Register new user
	cmds.register("reset", handleReset)                             // Delete all users
	cmds.register("users", handleGetUsers)                          // List all users
	cmds.register("agg", handleAgg)                                 // Scrape posts from feeds
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))     // Insert feed to table
	cmds.register("feeds", handleFeeds)                             // Gets all feeds and related user
	cmds.register("follow", middlewareLoggedIn(handleFollow))       // Create follow between user and feed
	cmds.register("following", middlewareLoggedIn(handleFollowing)) // List all followed feeds for current user
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))   // Deletes a followed relation between user and feed
	cmds.register("browse", middlewareLoggedIn(handleBrowse))       // Prints posts from followed feeds

	// Checks if user enters a command followed by one or
	// more arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	// Obtaining the command name and
	// arguments from text input
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Runs command using the input cmdName and cmdArgs
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
