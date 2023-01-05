package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mata649/mail_indexer/email"
	"github.com/mata649/mail_indexer/paths"
)

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

	emailPaths, err := paths.GetEmailsPaths(mainPath)
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

	email.GetEmails(emailPaths[:100000], currentDir)

	fmt.Println(time.Since(start))
}
