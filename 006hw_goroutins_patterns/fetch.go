package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func fetchFromAPI(wg *sync.WaitGroup, url string, ch chan<- []APIResponse) {
	defer wg.Done()
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request to %s: %v", url, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed request: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response from %s: %v\n", url, err)
		return
	}

	var apiResponses []APIResponse
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		fmt.Printf("Error while decoding JSON: %v\n", err)
		return
	}

	ch <- apiResponses
}

func aggregateResult(apiResponsesChannel <-chan []APIResponse, tasks chan<- DownloadTask) {
	for c := range apiResponsesChannel {
		for _, task := range c {
			tasks <- DownloadTask{Id: task.Id, DownloadUrl: task.DownloadUrl}
		}
	}
	close(tasks)
}
