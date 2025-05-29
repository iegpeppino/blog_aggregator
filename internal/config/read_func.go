package config

import (
	"encoding/json"
	"io"
	"os"
)

const gatorconfig_path = "/.gatorconfig.json"

// Reads the jsonfile with the configuration
// for the postgress server and returns
// a config struct
func Read() (Config, error) {

	file_path, err := getConfigPathFile()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(file_path)
	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	//fmt.Println("DB_URL: ", config.DB_URL)
	//fmt.Println("Current User Name: ", config.CurrentUserName)

	return config, nil
}

func getConfigPathFile() (string, error) {

	home_path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	file_path := home_path + gatorconfig_path

	return file_path, nil
}
