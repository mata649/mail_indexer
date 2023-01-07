package email

// import (
// 	"bufio"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"testing"
// )

// func TestAddLine(t *testing.T) {
// 	// create a temporary file
// 	tempFile, err := ioutil.TempFile("", "test-line-*")
// 	if err != nil {
// 		t.Fatalf("Failed to create temporary file: %v", err)
// 	}
// 	defer tempFile.Close()
// 	defer os.Remove(tempFile.Name())

// 	addLine("line1", tempFile)

// 	scanner := bufio.NewScanner(tempFile)
// 	scanner.Split(bufio.ScanLines)
// 	tempFile.Seek(0, 0)
// 	var contents string
// 	for scanner.Scan() {
// 		contents = scanner.Text()
// 	}

// 	// check that the line was added correctly
// 	expected := "line1"
// 	if contents != expected {
// 		t.Errorf("Expected file contents to be %q but got %q", expected, contents)
// 	}
// }

// func TestSaveEmails(t *testing.T) {

// 	testEmails := []Email{
// 		{
// 			MessageID: "test-1",
// 			Date:      "01/01/2022",
// 			From:      "test@example.com",
// 			To:        []string{"recipient1@example.com", "recipient2@example.com"},
// 			Subject:   "Test Email 1",
// 			Content:   "This is the content of test email 1",
// 		},
// 		{
// 			MessageID: "test-2",
// 			Date:      "01/02/2022",
// 			From:      "test@example.com",
// 			To:        []string{"recipient1@example.com"},
// 			Subject:   "Test Email 2",
// 			Content:   "This is the content of test email 2",
// 		},
// 	}

// 	tempDir := os.TempDir()
// 	tempFile := filepath.Join(tempDir, "file1.ndjson")
// 	saveEmails(testEmails, 1, tempDir)
// 	file, err := os.Open(tempFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	os.Remove(tempFile)
// 	defer file.Close()
// 	expectedLines := [4]string{
// 		`{ "index" : { "_index" : "emails" } }`,
// 		`{"messageID":"test-1","date":"01/01/2022","from":"test@example.com","to":["recipient1@example.com","recipient2@example.com"],"subject":"Test Email 1","content":"This is the content of test email 1"}`,
// 		`{ "index" : { "_index" : "emails" } }`,
// 		`{"messageID":"test-2","date":"01/02/2022","from":"test@example.com","to":["recipient1@example.com"],"subject":"Test Email 2","content":"This is the content of test email 2"}`,
// 	}
// 	counter := 0
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if expectedLines[counter] != line {
// 			t.Errorf("Expected line to be \n %v \n but got \n %v \n", expectedLines[counter], line)
// 		}
// 		counter += 1
// 	}

// }
