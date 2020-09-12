package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	cpuNumber := runtime.NumCPU()
	desiredThreads, err := strconv.Atoi(os.Args[1])
	check(err)
	threadNumber := desiredThreads
	if desiredThreads > cpuNumber {
		threadNumber = cpuNumber
	}
	url := os.Args[3]
	requests, err := strconv.Atoi(os.Args[2])
	check(err)

	logFile := os.Args[4]
	for i := 0; i < threadNumber; i++ {
		done := make(chan string)
		go executeThread(requests, url, logFile, done)
		println(<-done)
	}
}

func executeThread(requests int, url string, file string, done chan string) {
	for i := 0; i < requests; i++ {
		done := make(chan int64)
		go executeRequest(url, file, done)
		println(<-done)
	}
	done <- "Finished"
}

func executeRequest(url string, file string, done chan int64) {
	random := rand.Intn(50)
	var resp *http.Response
	var err error
	if random == 1 {
		startTime := time.Now()
		resp, err = http.Get(url)
		endTime := time.Now()
		durationTime := endTime.Sub(startTime).Microseconds()
		check(err)
		if resp.StatusCode == 200 {
			saveTextToFile(file, strconv.Itoa(int(durationTime)))
		}
	} else {
		resp, err = http.Get(url)
		check(err)
	}

	done <- 0
}

func saveTextToFile(file string, value string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	response, err := f.WriteString(value + "\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", response)

	syncError := f.Sync()
	check(syncError)
	closeError := f.Close()
	check(closeError)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
