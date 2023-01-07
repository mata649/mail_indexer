package email

import (
	"bytes"
	"encoding/json"
)

// // Appends a line to a file.
// // It takes the line to append and the file as arguments.
// func addLine(line string, memFile string) {

// 	memFile += line + "\n"

// }

// Converts a slice of Email structs to a slice of JSON strings and
// appends them to a file. The file is created if it does not exist.
func saveEmails(emails []Email) *bytes.Buffer {
	memFile := bytes.NewBufferString("")
	for _, email := range emails {
		parsedEmail, _ := json.Marshal(email)
		memFile.WriteString(`{ "index" : { "_index" : "emails" } }` + "\n")
		memFile.WriteString(string(parsedEmail) + "\n")

	}
	return memFile
}
