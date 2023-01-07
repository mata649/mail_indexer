package zinc

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/mata649/mail_indexer/pkg/config"
)

type ZincResponse struct {
	Message     string `json:"message"`
	RecordCount int    `json:"record_count"`
}

// Makes a POST request to the specified zinc host with the provided bytes buffer in the request body.
// It takes in a pointer to a bytes.Buffer bytesEmail and a pointer to a Configuration, returns a ZincResponse struct.
// The request is authenticated using basic auth with the provided username and password.
// The content type of the request is set to "application/json".
// The response body is read and unmarshalled into a ZincResponse struct.
// If any errors occur during the request or response processing, they are logged and the program exits.
func MakeRequest(bytesEmail *bytes.Buffer, currentConfig *config.Configuration) ZincResponse {

	req, err := http.NewRequest("POST", currentConfig.ZincHost+"/api/_bulk", bytesEmail)
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
