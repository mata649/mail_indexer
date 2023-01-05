package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	// Test loading configuration from a file
	configData := `{
		"NWorkers": 2,
		"EmailsPerFile": 10,
		"ZincHost": "localhost",
		"User": "user",
		"Password": "password"
	}`
	tmpFile, err := ioutil.TempFile("", "config-*.json")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	if _, err := tmpFile.Write([]byte(configData)); err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}
	config, err := LoadConfiguration(tmpFile.Name())
	os.Remove(tmpFile.Name())
	if err != nil {
		t.Errorf("Error loading configuration from file: %v", err)
	}
	if config.NWorkers != 2 {
		t.Errorf("Unexpected value for NWorkers: got %v, want %v", config.NWorkers, 2)
	}
	if config.EmailsPerFile != 10 {
		t.Errorf("Unexpected value for EmailsPerFile: got %v, want %v", config.EmailsPerFile, 10)
	}
	if config.ZincHost != "localhost" {
		t.Errorf("Unexpected value for ZincHost: got %v, want %v", config.ZincHost, "localhost")
	}
	if config.User != "user" {
		t.Errorf("Unexpected value for User: got %v, want %v", config.User, "user")
	}
	if config.Password != "password" {
		t.Errorf("Unexpected value for Password: got %v, want %v", config.Password, "password")
	}

	// Test loading configuration from a non-existent file
	_, err = LoadConfiguration("non-existent-file.json")
	if err == nil {
		t.Error("Expected error loading configuration from non-existent file, got nil")
	}

	// Test loading configuration from a file
	configData = `{
		"NWorker": 2,
		"EmailsPrFile": 10,
		"ZincHot": "localhost",
		"Use": "user",
	}`
	tmpFile, err = ioutil.TempFile("", "config-*.json")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	if _, err := tmpFile.Write([]byte(configData)); err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}
	_, err = LoadConfiguration(tmpFile.Name())
	if err == nil {
		t.Error("Expected error loading configuration from json with wrong structure, got nil")
	}
	os.Remove(tmpFile.Name())
}
