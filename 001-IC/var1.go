package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Message struct {
	Token  string
	FileID string
	Data   string
}

// TokenValidator interface for token validation
type TokenValidator interface {
	Validate(token string) bool
}

// MainTokenValidator validates tokens against a set of allowed tokens
type MainTokenValidator struct {
	allowedTokens map[string]bool
}

// NewMainTokenValidator creates a new instance of MainTokenValidator
func NewMainTokenValidator(tokens []string) *MainTokenValidator {
	allowedTokens := make(map[string]bool)
	for _, token := range tokens {
		allowedTokens[token] = true
	}
	return &MainTokenValidator{allowedTokens: allowedTokens}
}

// Validate checks if a token is valid
func (mtv *MainTokenValidator) Validate(token string) bool {
	return mtv.allowedTokens[token]
}

// MessageCache interface for message caching
type MessageCache interface {
	AddMessage(msg Message)
	FlushToFiles()
}

// MainMsgCache caches messages and flushes them to files
type MainMsgCache struct {
	mu        sync.RWMutex
	cache     map[string][]Message
	validator TokenValidator
}

// NewMainMsgCache creates a new instance of MainMsgCache
func NewMainMsgCache(validator TokenValidator) *MainMsgCache {
	return &MainMsgCache{
		cache:     make(map[string][]Message),
		validator: validator,
	}
}

// AddMessage adds a message to the cache
func (mmc *MainMsgCache) AddMessage(msg Message) {
	mmc.mu.Lock()
	defer mmc.mu.Unlock()

	if mmc.validator.Validate(msg.Token) {
		mmc.cache[msg.FileID] = append(mmc.cache[msg.FileID], msg)
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing file: ", err)
	}
}

func (mmc *MainMsgCache) RetryWrite2File(fileID string, messages []Message) {
	filePath := fmt.Sprintf("%s.txt", fileID)
	var success bool

	for retry := 0; retry < 5; retry++ {
		success = true

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("File not found or not created", err)
			success = false
			continue
		}

		for _, msg := range messages {
			_, err := file.WriteString(msg.Data + "\n")
			if err != nil {
				log.Printf("Failed to write to file %s: %v\n", filePath, err)
				success = false
				break
			}
		}
		closeFile(file)

		if success {
			break
		}

		log.Printf("Retrying to write to file %s (attemt %d)\n", filePath, retry+1)
		time.Sleep(10 * time.Second)
	}

	if success {
		mmc.mu.Lock()
		defer mmc.mu.Unlock()
		delete(mmc.cache, fileID)
	} else {
		log.Printf("Failed to write to file %s after retries\n", filePath)
	}
}

func (mmc *MainMsgCache) FlushToFiles() {
	mmc.mu.Lock()
	defer mmc.mu.Unlock()

	for fileID, message := range mmc.cache {
		go mmc.RetryWrite2File(fileID, message)
	}
}

// func (mmc *MainMsgCache) FlushToFiles() {
// 	mmc.mu.Lock()
// 	defer mmc.mu.Unlock()

// 	for fileID, messages := range mmc.cache {
// 		filePath := fmt.Sprintf("%s.txt", fileID)

// 		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
// 		if err != nil {
// 			fmt.Println("File not found or not created", err)
// 			continue
// 		}

// 		for _, msg := range messages {
// 			_, err := file.WriteString(msg.Data + "\n")
// 			if err != nil {
// 				log.Printf("Failed to write to file %s: %v\n", filePath, err)
// 			}
// 			continue
// 		}
// 		closeFile(file)

// 		delete(mmc.cache, fileID)
// 	}
// }

// Worker executes a task at regular intervals
type Worker struct {
	task     func()
	interval time.Duration
}

// NewWorker creates a new instance of Worker
func NewWorker(task func(), interval time.Duration) *Worker {
	return &Worker{
		task:     task,
		interval: interval,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context done, executing task before exit")
			w.task()
			log.Println("Stoping flush to file:", ctx.Err())
			return
		case <-ticker.C:
			log.Println("Ticker ticked, executing task")
			w.task()
		}
	}
}

// AddUsers simulates adding users and messages to message channels
func AddUsers(channels []chan Message, tokens []string, numMsg int) {
	for i := 1; i < numMsg; i++ {
		for user, ch := range channels {
			msg := Message{
				Token:  tokens[user],
				FileID: fmt.Sprintf("file_%d.txt", user),
				Data:   fmt.Sprintf("Message %d from user %d", i, user),
			}
			ch <- msg
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	// Config
	workerInterval := 2 * time.Second
	validTokens := []string{"token1", "token2", "token3"}
	numUsers := 3
	numMsg := 10
	countWorkers := 2

	validator := NewMainTokenValidator(validTokens)
	cache := NewMainMsgCache(validator)

	messageChannels := make([]chan Message, numUsers)
	for i := range messageChannels {
		messageChannels[i] = make(chan Message, 100)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := &sync.WaitGroup{}

	for i := 0; i < countWorkers; i++ {
		wg.Add(1)
		worker := NewWorker(cache.FlushToFiles, workerInterval)
		go worker.Start(ctx, wg)
	}

	for _, ch := range messageChannels {
		go func(mc chan Message) {
			for msg := range mc {
				cache.AddMessage(msg)
				fmt.Println("added to cache:", msg)
			}
		}(ch)
	}

	// Create users
	go AddUsers(messageChannels, validTokens, numMsg)

	<-ctx.Done()
	log.Println("Main: context done, waiting for worker to finish")

	wg.Wait()
	log.Println("Main: worker has been stopped")
}
