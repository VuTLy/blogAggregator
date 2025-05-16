package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/VuTLy/blogAggregator/internal/database"
	"github.com/google/uuid"
)

// handlerLogin sets the current user in the config file
func handlerLogin(s *state, cmd command) error {
	// Check if the command arguments are empty; if so, return an error indicating the username is required.
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required for login") // Return an error if no username is provided.
	}

	// Get the username from the first argument of the command
	username := cmd.args[0] // Set username to the first argument passed in the command.

	// Check if the user exists using GetUser
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		// If the error is not "sql: no rows in result set", log it
		if err.Error() != "sql: no rows in result set" {
			return fmt.Errorf("error checking user existence: %v", err)
		}

		// If the user doesn't exist, exit with code 1 and print an error
		fmt.Println("Error: user does not exist!")
		os.Exit(1) // Exit with code 1
	}

	// Set the current user in the configuration by calling the SetUser method on the config (s.cfg).
	err = s.cfg.SetUser(username) // Call SetUser method to update the configuration with the new username.
	if err != nil {               // If there is an error while setting the user (e.g., if writing to the config fails), return a formatted error.
		return fmt.Errorf("failed to set user: %v", err) // Return the error with a message indicating the failure.
	}

	// Print a success message to indicate the user has been successfully set in the config file.
	fmt.Printf("User %s has been set\n", username) // Print the success message with the username.
	return nil                                     // Return nil to indicate the operation was successful.
}

// handlerRegister creates a new user in the database
func handlerRegister(s *state, cmd command) error {
	// Ensure that the name was passed in the arguments
	if len(cmd.args) == 0 {
		return fmt.Errorf("name is required for registration")
	}

	// Use the provided name for the new user
	name := cmd.args[0]

	// Check if the user already exists using GetUser
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		// User already exists, return an error
		return fmt.Errorf("user '%s' already exists", name)
	}

	// If the error is anything other than "sql: no rows in result set", log it
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	// Generate a new UUID for the user
	userID := uuid.New()

	// Get the current time for created_at and updated_at
	now := time.Now()

	// Create the new user in the database
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Set the current user in the config
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to update config with new user: %v", err)
	}

	// Print success message
	fmt.Printf("User '%s' created successfully!\n", name)
	// Optionally log the user's details for debugging purposes
	fmt.Printf("Created user: %+v\n", newUser)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Println("Failed to reset database:", err)
		os.Exit(1)
	}
	fmt.Println("Database reset successfully.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	}

	// Sort users by name (case-insensitive)
	sort.Slice(users, func(i, j int) bool {
		return strings.ToLower(users[i].Name) < strings.ToLower(users[j].Name)
	})

	for _, user := range users {
		line := "* " + user.Name
		if user.Name == s.cfg.CurrentUserName {
			line += " (current)"
		}
		fmt.Println(line)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Link: %s\n\n", feed.Channel.Link)

	for _, item := range feed.Channel.Item {
		fmt.Printf("Title: %s\nLink: %s\nDate: %s\n\n", item.Title, item.Link, item.PubDate)
	}

	return nil
}

// Feeds portion
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	now := time.Now()
	feedID := uuid.New()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	// Automatically follow the feed
	followID := uuid.New()
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to auto-follow feed: %w", err)
	}

	fmt.Printf("Feed created and followed successfully: %+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeedsWithUser(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	fmt.Println("Feeds:")
	for _, feed := range feeds {
		fmt.Printf("* %s (%s) - added by %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("URL required to follow feed")
	}
	feedURL := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("could not find feed with URL: %s", feedURL)
	}

	now := time.Now()
	followID := uuid.New()

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("User '%s' is now following feed '%s'\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("could not get followed feeds: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Println("Feeds you are following:")
	for _, f := range follows {
		fmt.Printf("* %s\n", f.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: unfollow <feed-url>")
	}

	feedURL := cmd.args[0]

	err := s.db.DeleteFeedFollowByUserAndURL(context.Background(), database.DeleteFeedFollowByUserAndURLParams{
		UserID: user.ID,
		Url:    feedURL,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("You have unfollowed the feed: %s\n", feedURL)
	return nil
}
