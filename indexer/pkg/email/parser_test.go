package email

import (
	"os"
	"testing"
)

func TestParseEmailFromFile(t *testing.T) {
	// Create a temporary file and write sample email to it
	file, err := os.Create("/tmp/email.txt")
	if err != nil {
		t.Error("Failed to create temporary file:", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	_, err = file.WriteString("Message-ID: 123\nDate: Mon, 01 Jan 2018 15:04:05 +0000\nFrom: john@example.com\nTo: peter@example.com, lisa@example.com\nSubject: Test Email\nX-FileName: email.txt\n\nThis is a test email.\n")
	if err != nil {
		t.Error("Failed to write to temporary file:", err)
	}

	// Parse email from file
	file.Seek(0, 0)
	email := parseEmailFromFile(file)

	if email.MessageID != "123" {
		t.Errorf("Expected MessageID to be %v, got %v", "123", email.MessageID)
	}
	if email.Date != "Mon, 01 Jan 2018 15:04:05 +0000" {
		t.Errorf("Expected Date to be %v, got %v", "Mon, 01 Jan 2018 15:04:05 +0000", email.Date)
	}

	if email.From != "john@example.com" {
		t.Errorf("Expected From to be %v, got %v", "john@example.com", email.From)
	}
	if len(email.To) != 2 || email.To[0] != "peter@example.com" || email.To[1] != "lisa@example.com" {
		t.Error("Incorrect To:", email.To)
	}
	if email.Subject != "Test Email" {
		t.Errorf("Expected Subject to be %v, got %v", "Test Email", email.Subject)
	}
	if email.Content != "\nThis is a test email.\n" {
		t.Errorf("Expected Content to be %v, got %v", "This is a test email.", email.Content)
	}
}
