package email

import (
	"testing"
)

func TestCreateNdjsonBuffer(t *testing.T) {
	// Create test email struct
	email := Email{
		MessageID: "123456",
		Date:      "2022-01-01",
		From:      "sender@example.com",
		To:        []string{"receiver1@example.com", "receiver2@example.com"},
		Subject:   "Test Email",
		Content:   "This is a test email",
	}

	// Create slice of test emails
	emails := []Email{email, email}

	// Create expected NDJSON buffer
	expectedNdjson := `{ "index" : { "_index" : "emails" } }
{"messageID":"123456","date":"2022-01-01","from":"sender@example.com","to":["receiver1@example.com","receiver2@example.com"],"subject":"Test Email","content":"This is a test email"}
{ "index" : { "_index" : "emails" } }
{"messageID":"123456","date":"2022-01-01","from":"sender@example.com","to":["receiver1@example.com","receiver2@example.com"],"subject":"Test Email","content":"This is a test email"}` + "\n"

	resultNdjson := createNdjsonBuffer(emails)

	// Compare result with expected NDJSON
	resultString := resultNdjson.String()
	if resultString != expectedNdjson {
		t.Errorf("Expected NDJSON buffer to be \n %s \n got \n %s \n", expectedNdjson, resultNdjson.String())
	}
}
