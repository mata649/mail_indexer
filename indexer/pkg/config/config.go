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

var CurrentConfig *Configuration

// Loads the configuration from the specified file and stores it in the global
// CurrentConfig variable. If the CurrentConfig has already been initialized, it will return the
// existing configuration. It takes in a string configName which represents the name of the configuration file.
// If an error occurs while opening the file, or decoding
// the file, an empty Configuration struct and the error will be returned.
func LoadConfiguration(configPath string) (*Configuration, error) {
	if CurrentConfig == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return &Configuration{}, err
		}
		defer file.Close()
		err = json.NewDecoder(file).Decode(&CurrentConfig)
		if err != nil {
			return &Configuration{}, err
		}
	}

	return CurrentConfig, nil
}
