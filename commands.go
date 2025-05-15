package main

import "fmt"

// commands struct holds a map of command names to their respective handler functions.
type commands struct {
	handlers map[string]func(*state, command) error // A map that stores handlers (functions) keyed by command name.
}

// run executes a command with the provided state
// This method looks up the command handler and executes it if found.
func (c *commands) run(s *state, cmd command) error {
	// Look up the handler function associated with the command's name.
	handler, exists := c.handlers[cmd.name] // Retrieve the handler function for the command name.

	if !exists { // If the handler doesn't exist, return an error indicating the command was not found.
		return fmt.Errorf("command '%s' not found", cmd.name) // Return a formatted error with the command name.
	}

	// If the handler exists, execute it with the provided state and command.
	return handler(s, cmd) // Call the handler with the state and command, and return the result.
}

// register adds a new handler for a command
// This method registers a new handler function for a specific command name.
func (c *commands) register(name string, f func(*state, command) error) {
	// If handlers map is nil, initialize it.
	if c.handlers == nil { // Check if the handlers map is nil (uninitialized).
		c.handlers = make(map[string]func(*state, command) error) // Initialize the handlers map if it's nil.
	}

	// Register the handler function for the provided command name.
	c.handlers[name] = f // Store the handler function in the map with the command name as the key.
}
