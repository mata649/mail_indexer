package zinc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mata649/mail_indexer/pkg/config"
)

func TestMakeRequest(t *testing.T) {

	expectedResponse := ZincResponse{
		Message:     "Success",
		RecordCount: 2,
	}

	// Test data
	bytesEmail := bytes.NewBufferString(`
		{ "index" : { "_index" : "emails" } }
		{"message_id":"123","date":"2022-01-01","from":"john@example.com","to":["mary@example.com"],"subject":"Hello","content":"Hello, Mary!"}
		{ "index" : { "_index" : "emails" } }
		{"message_id":"456","date":"2022-01-02","from":"mary@example.com","to":["john@example.com"],"subject":"Re: Hello","content":"Hi, John!"}
	`)
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	currentConfig, err := config.LoadConfiguration(filepath.Join(workDir, "..", "..", "config.testing.json"))
	if err != nil {
		panic(err)
	}

	actualResponse := MakeRequest(bytesEmail, currentConfig)

	// Comparing the response record count with the expected record count
	if expectedResponse.RecordCount != actualResponse.RecordCount {
		t.Errorf("Error: record expected to be: %v, got %v", expectedResponse.RecordCount, actualResponse.RecordCount)
	}
}
