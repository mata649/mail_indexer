package email

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func addLine(currentFilePath string, line string) {

	f, err := os.OpenFile(currentFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = io.WriteString(f, line+"\n")
	if err != nil {
		fmt.Println(err)
	}
}
func saveEmails(emails []Email, counter int, currentDir string) string {

	currentFilePath := filepath.Join(currentDir, fmt.Sprintf("file%v.ndjson", counter))
	_, err := os.Create(currentFilePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, email := range emails {
		parsedEmail, _ := json.Marshal(email)
		addLine(currentFilePath, `{ "index" : { "_index" : "emails" } }`)
		addLine(currentFilePath, string(parsedEmail))

	}
	return currentFilePath
}
