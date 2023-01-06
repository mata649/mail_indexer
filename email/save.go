package email

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Appends a line to a file.
// It takes the line to append and the file as arguments.
func addLine(line string, f *os.File) {

	_, err := io.WriteString(f, line+"\n")
	if err != nil {
		fmt.Println(err)
	}
}

// Converts a slice of Email structs to a slice of JSON strings and
// appends them to a file. The file is created if it does not exist.
func saveEmails(emails []Email, counter int, currentDir string) {

	currentFilePath := filepath.Join(currentDir, fmt.Sprintf("file%v.ndjson", counter))
	_, err := os.Create(currentFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saving: %v \n", currentFilePath)
	f, err := os.OpenFile(currentFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for _, email := range emails {
		parsedEmail, _ := json.Marshal(email)
		addLine(`{ "index" : { "_index" : "emails" } }`, f)
		addLine(string(parsedEmail), f)

	}
	log.Printf("Saved: %v \n", currentFilePath)
}
