package zinc

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ZincResponse struct {
	Message     string `json:"message"`
	RecordCount int    `json:"record_count"`
}

func MakeRequest(bytesEmail *bytes.Buffer) ZincResponse {

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
