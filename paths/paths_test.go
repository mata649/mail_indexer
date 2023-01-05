package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetMainPath(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the current work directory")
	}
	testPath := filepath.Join(path, "testing")
	// Test with a valid directory path
	mainPath, err := GetMainPath(filepath.Join(testPath, "validDirectory"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedMainPath := filepath.Join(testPath, "validDirectory", "maildir")
	if mainPath != expectedMainPath {
		t.Errorf("expected %q, got %q", expectedMainPath, mainPath)
	}

	// Test with an invalid directory path
	_, err = GetMainPath("invalid_path")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	// Test with a valid directory path that does not contain a maildir directory
	_, err = GetMainPath(filepath.Join(path, "test", "validWithoutMaildir"))
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestGetFilePaths(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the current work directory")
	}
	mainPath := filepath.Join(path, "testing", "validDirectory", "maildir")

	// Test with a valid directory path
	filePaths, err := GetFilePaths(mainPath)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedFilePaths := []string{
		filepath.Join(mainPath, "1"),
		filepath.Join(mainPath, "2"),
		filepath.Join(mainPath, "p-allen", "1"),
		filepath.Join(mainPath, "p-allen", "2"),
		filepath.Join(mainPath, "p-allen", "testing", "1"),
		filepath.Join(mainPath, "p-allen", "testing", "2"),
		filepath.Join(mainPath, "p-allen", "testing", "3"),
	}

	if !equal(filePaths, expectedFilePaths) {
		t.Errorf("expected %v, got %v", expectedFilePaths, filePaths)
	}

	// Test with an invalid directory path
	_, err = GetFilePaths("invalid_path")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestDividePaths(t *testing.T) {
	paths := []string{"path1", "path2", "path3", "path4", "path5"}

	// Test with step=2
	dividedPaths := DividePaths(paths, 2)
	expectedDividedPaths := [][]string{{"path1", "path2"}, {"path3", "path4"}, {"path5"}}
	if !equal2D(dividedPaths, expectedDividedPaths) {
		t.Errorf("expected %v, got %v", expectedDividedPaths, dividedPaths)
	}

	// Test with step=3
	dividedPaths = DividePaths(paths, 3)
	expectedDividedPaths = [][]string{{"path1", "path2", "path3"}, {"path4", "path5"}}
	if !equal2D(dividedPaths, expectedDividedPaths) {
		t.Errorf("expected %v, got %v", expectedDividedPaths, dividedPaths)
	}
}

// Returns true if two slices contain the same elements in the same order
// else, returns false
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// equal2D returns true if the two slices contain the same slices in the same order, false otherwise.
func equal2D(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !equal(v, b[i]) {
			return false
		}
	}
	return true
}
