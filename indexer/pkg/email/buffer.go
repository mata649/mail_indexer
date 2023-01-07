package email

import (
	"bytes"
	"encoding/json"
)

// Converts a slice of Email structs into a buffer of bytes in NDJSON format.
// It takes in a slice of Email structs emails and returns a pointer to a bytes.Buffer.
// The buffer is in the correct NDJSON format to be saved using the zinc.MakeRequest function.
// Each Email struct in the slice is marshalled into a JSON object, and then the "index" field is added to
// each object. The resulting strings are then written to the buffer, with a newline character separating each object.
func createNdjsonBuffer(emails []Email) *bytes.Buffer {
	memFile := bytes.NewBufferString("")
	for _, email := range emails {
		parsedEmail, _ := json.Marshal(email)
		memFile.WriteString(`{ "index" : { "_index" : "emails" } }` + "\n")
		memFile.WriteString(string(parsedEmail) + "\n")

	}
	return memFile
}
