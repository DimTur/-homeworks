package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func downloadImg(ctx context.Context, task DownloadTask) error {
	req, err := http.NewRequestWithContext(ctx, "GET", task.DownloadUrl, nil)
	if err != nil {
		return fmt.Errorf("file downloading failed %s: %v", task.DownloadUrl, err)
	}

	response, err := http.DefaultClient.Do(req)
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
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		err := downloadImg(ctx, task)
		cancel()
		results <- err
	}
}
