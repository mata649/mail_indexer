package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
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
	dirPath := *emailPath

	mainPath, err := paths.GetMainPath(dirPath)
	if err != nil {
		log.Panicf("%v : %v", err, dirPath)
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
	err = loadToZinc(currentDir)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))

}
