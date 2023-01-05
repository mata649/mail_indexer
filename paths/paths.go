package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Returns the root path of the emails in the directory provided
// if the directory does not exist or the maildir directory can not
// be found in the provided directory, the function returns an error
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

// Returns a slice of strings containing the paths of
// all the files in the given directory and its subdirectories.
// If the directory cannot be accessed, the function returns an error.
func GetFilePaths(path string) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	var filePaths []string
	err := filepath.WalkDir(path, func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {

			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return filePaths, nil
}

// Returns the path to the "data" directory with the current date and time appended.
// If the current working directory cannot be accessed, the function returns an error.
func GetCurrentDataPath() (string, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("work directory could not be gotten")
	}
	return filepath.Join(workDir, "data", currentTime), nil

}

// Divides the given slice of strings into smaller slices of the specified size.
// The function returns a slice of slices of strings.
func DividePaths(paths []string, step int) [][]string {
	var pathsDivided [][]string
	startIndex := 0
	endIndex := step
	slices := len(paths) / step
	totalPaths := len(paths)

	for i := 0; i < slices+1; i++ {
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
