package zinc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/mata649/mail_indexer/config"
)

func MakeRequest(filePath string, wg *sync.WaitGroup, semaphore chan bool, currentConfig config.Configuration) {
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
	fmt.Println(string(body))
	wg.Done()
	<-semaphore
}
