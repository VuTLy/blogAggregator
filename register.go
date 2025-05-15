package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/VuTLy/blogAggregator/internal/database"

	"github.com/google/uuid"
)

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
