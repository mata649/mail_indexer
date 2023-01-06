package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	NWorkers      int
	EmailsPerFile int
	ZincHost      string
	User          string
	Password      string
}

// Loads the configuration setted in the config.json file.
// If the config.json file can't be opened, or the
// config.json file doesn't have the correct structure
// the function will return an err
func LoadConfiguration(path string) (Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return Configuration{}, err
	}
	return configuration, nil
}
