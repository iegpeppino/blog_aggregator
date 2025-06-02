package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iegpeppino/blog_aggregator/internal/database"
)

// Scrapes posts from followed feeds and saves them
// to the posts table
// A time_between_reqs arg is passed to create a ticker
// so the scraping is done periodically using that time.Duration
func handleAgg(s *state, cmd command) error {
	// Validate argument
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("invalid argument, need <time_between_reqs> ")
	}

	// Parse the argument and create ticker
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration %w", err)
	}

	log.Printf("Fetching feeds every %s ...", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)

	// Use the ticker to call scrapeFeeds periodically
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

// Gets the next feed to fetch, marks it as fetched
// and then scrapes all posts from it, saving them to db
func scrapeFeeds(s *state) {

	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Println("Couldn't get next feed to fetch", err)
		return
	}

	_, err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched %v", feed.Name, err)
		return
	}

	feedContent, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("Couldn't parse feed %s %v", feed.Name, err)
		return
	}

	for _, item := range feedContent.Channel.Item {

		// Preparing to save post to db
		publishedAt := sql.NullTime{} // In case there's no PubDate data
		if time, err := time.Parse(time.RFC1123Z, item.PubDate); err != nil {
			publishedAt = sql.NullTime{
				Time:  time,
				Valid: true,
			}
		}

		// Setting new post parameters
		postParams := database.CreatePostParams{
			ID:    uuid.New(),
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{ // In case there's no description
				String: item.Description,
				Valid:  true,
			},
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}

		// Query the insert post into db
		_, err = s.db.CreatePost(ctx, postParams)
		if err != nil {
			// Ignoring existing url error
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post %v", err)
			continue
		}
	}
	log.Printf("Feed %s fetched, %v posts where saved", feed.Name, len(feedContent.Channel.Item))
}
