package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

type Source struct {
	Source Email `json:"_source"`
}
type Hits struct {
	Hits []Source `json:"hits"`
}
type ZincResponse struct {
	Hits Hits `json:"hits"`
}
type Email struct {
	MessageID string   `json:"messageID"`
	Date      string   `json:"date"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	Content   string   `json:"content"`
}

func searchInZincEngine(w http.ResponseWriter, r *http.Request) {

	text := r.URL.Query().Get("text")
	query := fmt.Sprintf(`{
	    "search_type": "match",
	    "query":
	    {
	        "term": "%v"
	    },
	    "from": 0,
	    "max_results": 20,
	    "_source": []
	}`, text)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/api/emails/_search", *zincHost), strings.NewReader(query))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(*zincUser, *zincPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	zincResponse := ZincResponse{}
	err = json.Unmarshal(body, &zincResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var emails []Email
	for _, source := range zincResponse.Hits.Hits {
		emails = append(emails, source.Source)
	}
	if len(emails) > 0 {

		json.NewEncoder(w).Encode(emails)
		return
	}
	json.NewEncoder(w).Encode([]string{})

}

func checkFlags() error {
	if *zincHost == "" {
		return fmt.Errorf("Zinc host is empty")
	}
	if *zincUser == "" {
		return fmt.Errorf("Zinc host is empty")
	}
	if *zincPassword == "" {
		return fmt.Errorf("Zinc host is empty")
	}
	return nil
}

var zincHost = flag.String("host", "", "zinc host")
var zincUser = flag.String("user", "", "zinc user")
var zincPassword = flag.String("password", "", "zinc password")

func main() {
	flag.Parse()
	err := checkFlags()
	if err != nil {
		panic(err)
	}
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		currentDir, err := os.Getwd()
		if err != nil {
			http.Error(w, "Error getting cwd", http.StatusInternalServerError)
		}
		staticDir := filepath.Join(currentDir, "static")
		if _, err := os.Stat(staticDir + r.URL.Path); errors.Is(err, os.ErrNotExist) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		}
		http.ServeFile(w, r, staticDir+r.URL.Path)
	})
	r.Get("/search", searchInZincEngine)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
