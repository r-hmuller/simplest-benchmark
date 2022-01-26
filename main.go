package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Jwt struct {
	Jwt string
}

func main() {
	cpuNumber := runtime.NumCPU()
	desiredThreads, err := strconv.Atoi(os.Args[1])
	check(err)
	threadNumber := desiredThreads
	if desiredThreads > cpuNumber {
		threadNumber = cpuNumber
	}
	baseUrl := os.Args[3]
	requests, err := strconv.Atoi(os.Args[2])
	check(err)

	print("Starting experiment with " + string(rune(threadNumber)) + " threads")

	var jsonStr = []byte(`{"username":"root", "password": "root"}`)
	req, err := http.NewRequest("POST", baseUrl+"/_open/auth", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	jwt := Jwt{}
	err = json.Unmarshal(body, &jwt)
	if err != nil {
		panic(err)
	}

	logFile := os.Args[4]
	for i := 0; i < threadNumber; i++ {
		done := make(chan string)
		go executeThread(requests, baseUrl, logFile, jwt, done)
		println(<-done)
	}
}

func executeThread(requests int, url string, file string, jwt Jwt, done chan string) {
	for i := 0; i < requests; i++ {
		done := make(chan int64)
		go executeRequest(url, file, jwt, done)
		println(<-done)
	}
	done <- "Finished"
}

func executeRequest(url string, file string, jwt Jwt, done chan int64) {
	random := rand.Intn(50)
	if random == 1 {
		startTime := time.Now()
		doRequest(url, jwt)
		endTime := time.Now()
		durationTime := endTime.Sub(startTime).Microseconds()
		saveTextToFile(file, strconv.Itoa(int(durationTime)))
	} else {
		doRequest(url, jwt)
	}

	done <- 0
}

func doRequest(url string, jwt Jwt) {
	var jsonStr = []byte(`{"username":"root", "password": "root"}`)
	req, err := http.NewRequest("POST", url+"/_open/auth", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "Bearer "+jwt.Jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(err)
	}
}

func saveTextToFile(file string, value string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	response, err := f.WriteString(strconv.FormatInt(time.Now().Unix(), 10) + "," + value + "\n")
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
