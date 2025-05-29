package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func (config *Config) SetUser(username string) error {

	config.CurrentUserName = username

	err := Write(config)
	if err != nil {
		return err
	}
	return nil
}

// I don't use a pointer as an argument because
// I dont need to modify the config struct
func Write(config *Config) error {

	// Marshal the Config struct into JSON
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error while Marshalin Config struct:")
		return err
	}

	// Get the filepath
	filepath, err := getConfigPathFile()
	if err != nil {
		fmt.Println("Error getting the filepath:")
		return err
	}

	// Opens file on over-write mode
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening the file:")
		return err
	}

	// Ensures the file is closed at the end
	defer file.Close()

	// Writes jsonConfig contents to the file
	_, err = file.Write(jsonConfig)
	if err != nil {
		fmt.Println("Error writing to file:")
		return err
	}

	fmt.Printf("JSON succesfully written to %v\n", filepath)
	return nil
}
