package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mata649/mail_indexer/pkg/config"
	"github.com/mata649/mail_indexer/pkg/email"
	"github.com/mata649/mail_indexer/pkg/paths"
	"github.com/mata649/mail_indexer/pkg/zinc"
)

var currentConfig config.Configuration

func saveEmails(emailPaths []string, currentDir string) {
	step := currentConfig.EmailsPerFile
	emailPathsDivided := paths.DividePaths(emailPaths, step)
	semaphore := make(chan bool, currentConfig.NWorkers)
	var wg sync.WaitGroup
	counter := 0
	for _, emailSlice := range emailPathsDivided {
		wg.Add(1)
		counter += 1
		go email.GetEmails(emailSlice, currentDir, &wg, counter, semaphore)

	}

	wg.Wait()
}

func loadToZinc(currentDir string) error {
	filePaths, err := paths.GetFilePaths(currentDir)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	semaphore := make(chan bool, currentConfig.NWorkers)

	for _, filePath := range filePaths {
		wg.Add(1)
		go zinc.MakeRequest(filePath, &wg, semaphore, currentConfig)
	}
	wg.Wait()
	return nil
}
func main() {
	start := time.Now()

	if len(os.Args) == 1 {
		panic("a directory has to be provided")
	}
	dirPath := os.Args[1]

	mainPath, err := paths.GetMainPath(dirPath)
	if err != nil {
		panic(err)
	}

	emailPaths, err := paths.GetFilePaths(mainPath)
	if err != nil {
		panic(err)
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	workDir, err := os.Getwd()
	if err != nil {
		panic("Work directory could not be gotten")
	}
	currentConfig, err = config.LoadConfiguration(filepath.Join(workDir, "config.json"))
	if err != nil {
		panic(err)
	}
	currentDir := filepath.Join(workDir, "data", currentTime)
	os.MkdirAll(currentDir, 0777)
	saveEmails(emailPaths, currentDir)
	// err = loadToZinc(currentDir)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println(time.Since(start))
}
