package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func downloadImg(task DownloadTask) error {
	response, err := http.Get(task.DownloadUrl)
	if err != nil {
		return fmt.Errorf("file downloading failed %s: %v", task.DownloadUrl, err)
	}
	defer response.Body.Close()

	out, err := os.Create(fmt.Sprintf("%s/%s", outputDir, task.Id))
	if err != nil {
		return fmt.Errorf("file creation failed %s: %v", task.Id, err)
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write %v to file %s: %v", task.DownloadUrl, task.Id, err)
	}
	return nil
}

func worker(wg *sync.WaitGroup, tasks <-chan DownloadTask, results chan<- error) {
	defer wg.Done()
	for task := range tasks {
		err := downloadImg(task)
		results <- err
	}
}
