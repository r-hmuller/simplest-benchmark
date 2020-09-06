package main

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {

	cpuNumber := runtime.NumCPU()

	url := os.Args[1]
	requests, error := strconv.Atoi(os.Args[2])
	if error != nil {
		print(error)
		os.Exit(1)
	}
	for i := 0; i < cpuNumber; i++ {
		done := make(chan string)
		go executeThread(requests, url, done)
		print(<-done)
	}
}

func executeThread(requests int, url string, done chan string) {
	println(requests)
	for i := 0; i < requests; i++ {
		done := make(chan int64)
		go executeRequest(url, done)
		println(<-done)
	}
	done <- "Finished"
}

func executeRequest(url string, done chan int64) {
	startTime := time.Now()
	resp, _ := http.Get(url)
	if resp.StatusCode == 200 {
		print("Ok")
	}
	endTime := time.Now()
	durationTime := endTime.Sub(startTime).Nanoseconds()
	done <- durationTime
}
