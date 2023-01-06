package zinc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/mata649/mail_indexer/pkg/config"
)

type ZincResponse struct {
	Message     string `json:"message"`
	RecordCount int    `json:"record_count"`
}

// Makes a request to the given URL with the provided configuration.
// The request is a POST request with the JSON payload read from the specified file path.
// If the request is successful, the function prints the status code and the response body.
// Otherwise, it prints an error message.
func MakeRequest(filePath string, wg *sync.WaitGroup, semaphore chan bool, currentConfig config.Configuration) ZincResponse {
	semaphore <- true
	jsonStr, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", currentConfig.ZincHost+"/api/_bulk", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(currentConfig.User, currentConfig.Password)
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
	wg.Done()
	<-semaphore
	return zincResponse
}
