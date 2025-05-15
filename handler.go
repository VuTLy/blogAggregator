package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
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
