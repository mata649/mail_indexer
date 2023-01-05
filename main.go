package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mata649/mail_indexer/email"
	"github.com/mata649/mail_indexer/paths"
)

func dividePaths(paths []string, step int) [][]string {
	var pathsDivided [][]string
	startIndex := 0
	endIndex := step
	slices := len(paths) / step
	totalPaths := len(paths)
	for i := 0; i < slices; i++ {
		if endIndex > len(paths) {
			pathsDivided = append(pathsDivided, paths[startIndex:totalPaths])
			break
		} else {
			pathsDivided = append(pathsDivided, paths[startIndex:endIndex])

		}
		startIndex = endIndex
		endIndex += step
	}
	return pathsDivided
}
func saveEmails(emailPaths []string, currentDir string) {
	step := 1000
	emailPathsDivided := dividePaths(emailPaths, step)
	semaphore := make(chan bool, 10)
	var wg sync.WaitGroup
	counter := 0
	for _, emailSlice := range emailPathsDivided {
		wg.Add(1)
		counter += 1
		go email.GetEmails(emailSlice, currentDir, &wg, counter, semaphore)

	}

	wg.Wait()
}

// func loadToZinc(currentDir string) error {
// 	filePaths, err := paths.GetFilePaths(currentDir)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
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
	currentDir := filepath.Join(workDir, "data", currentTime)
	os.MkdirAll(currentDir, 0777)
	saveEmails(emailPaths, currentDir)
	// loadToZinc(currentDir)
	fmt.Println(time.Since(start))
}
