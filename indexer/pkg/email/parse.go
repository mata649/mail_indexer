package email

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Reads the given email paths and stores the data in Email structs.
// Then, the emails are saved in files in the currentDir directory.
// The counter variable is used to number the files.
// A semaphore channel is passed to the function to control the concurrent access to the function, limited by nWorkers.
// A wait group is passed to the function to control the goroutines.
func ParseEmails(emailPaths []string, currentDir string, wg *sync.WaitGroup, counter int, semaphore chan bool) {
	semaphore <- true
	var emails []Email

	for _, path := range emailPaths {

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var email Email
		currentParam := ""
		mainParamsChecked := false
		for scanner.Scan() {
			line := scanner.Text()
			firstLineParam := false
			var lineSplitted []string
			if !mainParamsChecked {

				lineSplitted = strings.Split(line, ":")
			}
			if len(lineSplitted) > 1 {
				currentParam = lineSplitted[0]
				if currentParam == "X-FileName" {

					mainParamsChecked = true
					continue
				}
				firstLineParam = true
			}
			if mainParamsChecked == false {
				switch currentParam {
				case "Message-ID":
					email.MessageID = strings.Trim(lineSplitted[1], " ")
				case "Date":
					email.Date = strings.Trim(strings.Join(lineSplitted[1:], ":"), " ")
				case "From":
					email.From = lineSplitted[1]
				case "To":
					if firstLineParam {
						emailsSplited := strings.Split(lineSplitted[1], ",")
						if len(emailsSplited) < 2 {
							email.To = append(email.To, strings.Trim(lineSplitted[1], " "))
							continue
						}
						for _, emailSplited := range emailsSplited {
							email.To = append(email.To, strings.Trim(emailSplited, " "))
						}
						continue
					}
					for _, emailSplited := range strings.Split(line, ",") {
						email.To = append(email.To, strings.Trim(strings.Trim(emailSplited, "	"), " "))

					}
				case "Subject":
					if firstLineParam {

						email.Subject = strings.Trim(strings.Join(lineSplitted[1:], ":"), " ")
						continue
					}
					email.Subject += lineSplitted[0]

				}
			} else if mainParamsChecked {
				email.Content += line + "\n"
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}
		}
		emails = append(emails, email)

	}
	bytesEmail := saveEmails(emails)
	resp := MakeRequestV2(bytesEmail)
	fmt.Println(resp)
	wg.Done()
	<-semaphore
}

type ZincResponse struct {
	Message     string `json:"message"`
	RecordCount int    `json:"record_count"`
}

func MakeRequestV2(bytesEmail *bytes.Buffer) ZincResponse {

	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulk", bytesEmail)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

	zincResponse := ZincResponse{}
	err = json.Unmarshal(body, &zincResponse)
	if err != nil {
		log.Fatal(err)
	}
	return zincResponse
}
