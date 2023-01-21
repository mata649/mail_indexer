package email

import (
	"fmt"
	"os"
	"sync"

	"github.com/mata649/mail_indexer/pkg/config"
	"github.com/mata649/mail_indexer/pkg/zinc"
)

// Makes a request to the zinc host to ingest a batch of emails stored in the given paths.
// Takes in a slice of strings emailPaths, a pointer to a sync.WaitGroup wg, a channel semaphore,
// and a pointer to a Configuration struct currentConfig.
// Iterates through the emailPaths slice, opening each file and parsing it into an Email struct using the
// getEmailFromFile function. Appends each Email struct to a slice of emails.
// Creates a buffer of bytes in NDJSON format using the createNdjsonBuffer function and the emails slice.
// Makes a request to the zinc host using the MakeRequest function, passing the bytesEmail buffer and the
// currentConfig struct.
func MakeIngestion(emailPaths []string, wg *sync.WaitGroup, semaphore chan bool, currentConfig *config.Configuration) {
	defer wg.Done()
	semaphore <- true
	defer func() {
		<-semaphore
	}()

	var emails []Email

	for _, path := range emailPaths {

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		email := parseEmailFromFile(file)
		emails = append(emails, email)

	}

	bytesEmail := createNdjsonBuffer(emails)
	resp := zinc.MakeRequest(bytesEmail, currentConfig)
	fmt.Println(resp)

}
