package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetMainPath(DirPath string) (string, error) {
	_, err := os.ReadDir(DirPath)
	if err != nil {
		return "", fmt.Errorf("the directory provided does not exist")
	}
	emailsDirPath := filepath.Join(DirPath, "maildir")
	_, err = os.ReadDir(emailsDirPath)
	if err != nil {
		return "", fmt.Errorf("the maildir directory was not found into the provided Dir")
	}
	return emailsDirPath, nil
}

func GetFilePaths(path string) ([]string, error) {
	var emailsPaths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			emailsPaths = append(emailsPaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return emailsPaths, nil
}
func GetCurrentDataPath() (string, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("work directory could not be gotten")
	}
	return filepath.Join(workDir, "data", currentTime), nil

}
