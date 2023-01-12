package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/mata649/mail_indexer/pkg/config"
	"github.com/mata649/mail_indexer/pkg/email"
	"github.com/mata649/mail_indexer/pkg/paths"
)

// saveEmails processes a slice of email paths in parallel using a number of workers specified in the configuration file.
// It takes in a slice of strings emailPaths which represents the paths of the emails to be processed.
// The emails are divided into smaller slices and sent to the email.MakeIngestion function in separate goroutines.
// A semaphore is used to control the number of concurrent goroutines.
// The function blocks until all goroutines have finished processing the emails.
func saveEmails(emailPaths []string) {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	currentConfig, err := config.LoadConfiguration(filepath.Join(workDir, "config.json"))
	if err != nil {
		panic(err)
	}

	step := currentConfig.EmailsPerFile
	emailPathsDivided := paths.DividePaths(emailPaths, step)
	semaphore := make(chan bool, currentConfig.NWorkers)
	defer close(semaphore)
	var wg sync.WaitGroup
	for _, emailSlice := range emailPathsDivided {
		wg.Add(1)
		go email.MakeIngestion(emailSlice, &wg, semaphore, currentConfig)

	}

	wg.Wait()
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var emailPath = flag.String("emailpath", "", "email location")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	start := time.Now()

	if *emailPath == "" {
		panic("a directory has to be provided")
	}

	mainPath, err := paths.GetMainPath(*emailPath)
	if err != nil {
		log.Panicf("%v : %v", err, *&mainPath)
	}

	emailPaths, err := paths.GetFilePaths(mainPath)
	if err != nil {
		panic(err)
	}

	saveEmails(emailPaths)
	fmt.Println(time.Since(start))

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
