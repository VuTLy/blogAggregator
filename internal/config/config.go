package config

import (
	"encoding/json" // Import the encoding/json package to handle JSON serialization and deserialization.
	"os"            // Import the os package to interact with the file system (open files, read environment variables, etc.).
	"path/filepath" // Import the filepath package to work with file paths in a cross-platform way.
)

const configFileName = ".gatorconfig.json" // Define a constant for the configuration file name.

type Config struct { // Define a struct named Config to hold configuration data.
	DBURL           string `json:"db_url"`            // Field to store the database URL, serialized as "db_url" in JSON.
	CurrentUserName string `json:"current_user_name"` // Field to store the current user's name, serialized as "current_user_name" in JSON.
}

// SetUser method to update the CurrentUserName field of the Config struct and save the updated config to a file.
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName // Set the CurrentUserName field to the provided userName.
	return write(*cfg)             // Call the write function to save the updated config to a file.
}

// Read function to read the configuration from the file and return it as a Config struct.
func Read() (Config, error) {
	fullPath, err := getConfigFilePath() // Get the full path to the config file.
	if err != nil {                      // If there is an error getting the file path, return an empty Config and the error.
		return Config{}, err
	}

	file, err := os.Open(fullPath) // Open the config file for reading.
	if err != nil {                // If there is an error opening the file, return an empty Config and the error.
		if os.IsNotExist(err) {
			// File doesn't exist, create a default config
			defaultConfig := Config{
				DBURL:           "", // Or some default URL
				CurrentUserName: "", // No user by default
			}

			// Write the default config to disk
			if err := write(defaultConfig); err != nil {
				return Config{}, err
			}

			return defaultConfig, nil
		}
		// If it's some other error, return it
		return Config{}, err
	}
	defer file.Close() // Ensure the file is closed when the function finishes.

	decoder := json.NewDecoder(file) // Create a new JSON decoder to read the JSON data from the file.
	cfg := Config{}                  // Create an empty Config struct to store the decoded data.
	err = decoder.Decode(&cfg)       // Decode the JSON data into the cfg variable.
	if err != nil {                  // If there is an error decoding the JSON, return an empty Config and the error.
		return Config{}, err
	}

	return cfg, nil // Return the populated Config and no error.
}

// getConfigFilePath function to construct the full path of the config file based on the user's home directory.
func getConfigFilePath() (string, error) {
	// Get the home directory
	dir, err := os.UserHomeDir()
	if err != nil { // If there is an error getting the home directory, return an empty string and the error.
		return "", err
	}
	fullPath := filepath.Join(dir, configFileName) // Construct the full file path by joining the home directory with the config file name.
	return fullPath, nil                           // Return the full file path and no error.
}

// write function to write the provided Config struct to the config file in JSON format.
func write(cfg Config) error {
	fullPath, err := getConfigFilePath() // Get the full path to the config file.
	if err != nil {                      // If there is an error getting the file path, return the error.
		return err
	}

	file, err := os.Create(fullPath) // Create or truncate the config file for writing.
	if err != nil {                  // If there is an error creating the file, return the error.
		return err
	}
	defer file.Close() // Ensure the file is closed when the function finishes.

	encoder := json.NewEncoder(file) // Create a new JSON encoder to write the config data to the file.
	err = encoder.Encode(cfg)        // Encode the Config struct and write it to the file.
	if err != nil {                  // If there is an error encoding the data, return the error.
		return err
	}

	return nil // Return nil if no errors occurred (write operation was successful).
}
