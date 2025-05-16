package main

import (
	"context"
	"fmt"

	"github.com/VuTLy/blogAggregator/internal/database" // adjust if needed
)

// middlewareLoggedIn wraps a handler that requires a logged-in user.
// It fetches the user based on the current config and passes it to the handler.
func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		username := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("must be logged in: %w", err)
		}
		return handler(s, cmd, user)
	}
}
