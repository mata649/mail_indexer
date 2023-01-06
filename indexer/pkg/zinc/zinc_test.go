package zinc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/mata649/mail_indexer/pkg/config"
)

func TestMakeRequest(t *testing.T) {
	currentWorkDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the current work directory")
	}

	currentConfig, err := config.LoadConfiguration(filepath.Join(currentWorkDir, "..", "config.testing.json"))
	if err != nil {
		t.Fatalf("Error getting the testing configuration")
	}

	// Create a test file
	tempFile, err := ioutil.TempFile("", "test_file")
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name()) // clean up

	// Write some JSON data to the file
	jsonData := `{ "index" : { "_index" : "emails" } }
	{"messageID":"test-1","date":"01/01/2022","from":"test@example.com","to":["recipient1@example.com","recipient2@example.com"],"subject":"Test Email 1","content":"This is the content of test email 1"}`
	if _, err := tempFile.Write([]byte(jsonData)); err != nil {
		t.Errorf("Error writing to temp file: %v", err)
	}

	// Initialize a wait group and semaphore channel for the test
	wg := sync.WaitGroup{}
	wg.Add(1)
	semaphore := make(chan bool, 1)

	// Call the MakeRequest function
	zincReponse := MakeRequest(tempFile.Name(), &wg, semaphore, currentConfig)
	wg.Wait()

	// Check the response for correctness
	// TODO: write the remainder of the test

	if zincReponse.RecordCount == 0 {
		t.Errorf("The email was not found in the zinc engine")
	}
}
