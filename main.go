package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/VuTLy/blogAggregator/internal/config"
	"github.com/VuTLy/blogAggregator/internal/database"
)

func main() {
	// Check if enough command-line arguments are provided. The program expects at least one argument (the command).
	if len(os.Args) < 2 {
		log.Fatal("Error: Not enough arguments. A command is required.") // If not enough arguments, print an error and exit.
	}

	// Initialize the state by reading the configuration file.
	cfg, err := config.Read() // Read the configuration from the config file.
	if err != nil {           // If an error occurs while reading the config, log the error and exit.
		log.Fatalf("Error reading config: %v", err) // Log the error and exit the program.
	}

	// Step 7: Open a connection to the database
	// Make sure your config struct has a field for the DB URL
	dbURL := cfg.DBURL // Adjust this based on your actual config structure
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close() // Ensure we close the database connection when the program exits

	// Create a new database queries instance
	dbQueries := database.New(db)

	// Create state with both config and database
	state := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	// Create a new commands struct and register the login handler.
	cmds := &commands{}                  // Initialize a new commands struct that will hold the handlers.
	cmds.register("login", handlerLogin) // Register the "login" command with the handlerLogin function.
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	// Parse the command and arguments from the command-line input.
	cmd := command{ // Create a new command struct.
		name: os.Args[1],  // Set the command name to the first argument (os.Args[1]).
		args: os.Args[2:], // Set the arguments to the rest of the command-line arguments (os.Args[2:] onwards).
	}

	// Run the command by invoking the appropriate handler for the command.
	err = cmds.run(state, cmd) // Call the run method of the cmds struct with the parsed state and command.
	if err != nil {            // If the command execution returns an error, log the error and exit.
		log.Fatal(err) // Log the error and terminate the program.
	}
}
