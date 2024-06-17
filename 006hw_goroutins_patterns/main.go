package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	workerCount = 8
	outputDir   = "./downloads"
)

func main() {
	var fetchWg sync.WaitGroup
	var workerWg sync.WaitGroup

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create output directory: %v\n", err)
			return
		}
	}

	apiUrls := []string{
		"https://picsum.photos/v2/list?page=2&limit=100",
		"https://picsum.photos/v2/list?page=3&limit=100",
	}

	// This is only true if limit=100 for each URL, otherwise we will have to write new logic..
	totalTasks := len(apiUrls) * 100

	apiResponsesChannel := make(chan []APIResponse, len(apiUrls))
	tasksChannel := make(chan DownloadTask, totalTasks)
	resultsChannel := make(chan error, totalTasks)

	startTime := time.Now()

	// Fetch data from API.
	for _, url := range apiUrls {
		fetchWg.Add(1)
		go fetchFromAPI(&fetchWg, url, apiResponsesChannel)
	}

	go func() {
		fetchWg.Wait()
		close(apiResponsesChannel)
	}()

	// Aggregate results and create download tasks.
	go aggregateResult(apiResponsesChannel, tasksChannel)

	// Start workers to download images.
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go worker(&workerWg, tasksChannel, resultsChannel)
	}

	// Wait for workers to complete and close results channel.
	go func() {
		workerWg.Wait()
		close(resultsChannel)
	}()

	// Process results.
	for err := range resultsChannel {
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("Download succeeded")
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("All downloads completed in %s\n", elapsedTime)
}
