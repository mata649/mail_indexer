package email

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/mata649/mail_indexer/config"
	"github.com/mata649/mail_indexer/paths"
)

// Please fill the emails directory with the emails to test
func TestGetEmails(t *testing.T) {
	// Previous params to test the function

	// Temp dir to save the data
	tempDataDir, err := ioutil.TempDir("", "data-*")
	if err != nil {
		t.Fatalf("Error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDataDir)

	currentWorkDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the current work directory")
	}
	testingDir := filepath.Join(currentWorkDir, "..", "testing")
	emailPaths, err := paths.GetFilePaths(filepath.Join(testingDir, "emails"))

	if err != nil {
		t.Fatalf("Error getting the emails")
	}
	currentConfig, err := config.LoadConfiguration(filepath.Join(currentWorkDir, "..", "config.testing.json"))
	if err != nil {
		t.Fatalf("Error getting the testing configuration")
	}

	step := currentConfig.EmailsPerFile
	emailPathsDivided := paths.DividePaths(emailPaths, step)

	// Making the directory to save the emails
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	dataDir := filepath.Join(tempDataDir, currentTime)
	os.MkdirAll(dataDir, 0777)

	semaphore := make(chan bool, currentConfig.NWorkers)
	var wg sync.WaitGroup
	counter := 0
	for _, emailSlice := range emailPathsDivided {
		wg.Add(1)
		counter += 1
		go GetEmails(emailSlice, dataDir, &wg, counter, semaphore)

	}
	wg.Wait()
	// Checking is correct
	files, _ := ioutil.ReadDir(dataDir)
	if len(files) != len(emailPathsDivided) {
		t.Errorf("Number of files is incorrect, got %v, want %v", len(files), len(emailPathsDivided))
	}
}
