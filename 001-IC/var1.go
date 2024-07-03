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
	WriteMsgs2Cache(msg Message)
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

// WriteMsgs2Cache adds a message to the cache
func (mmc *MainMsgCache) WriteMsgs2Cache(msg Message) {
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

// RetryWrite2File retries writing messages to a file with exponential backoff in case of failures.
// It attempts multiple times and returns an error if unsuccessful after the maximum retries.
func (mmc *MainMsgCache) RetryWrite2File(fileID string, messages []Message) error {
	const maxRetries = 5
	const baseDelay = 3 * time.Second

	filePath := fmt.Sprintf("%s.txt", fileID)
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			lastErr = err
			log.Printf("Attempt %d: failed to open or create file %s: %v", attempt, filePath, err)
		}
		defer closeFile(file)

		for _, msg := range messages {
			_, err := file.WriteString(msg.Data + "\n")
			if err != nil {
				lastErr = err
				log.Printf("Attempt %d: failed to write to file %s: %v", attempt, filePath, err)
				break
			}
		}
		if lastErr == nil {
			return nil
		}

		time.Sleep(baseDelay * time.Duration(1<<uint(attempt-1)))
	}

	return fmt.Errorf("failed to write to file %s after %d attempts: %v", filePath, maxRetries, lastErr)
}

// FlushToFiles flushes cached messages to respective files
// It removes successfully flushed entries from the cache
func (mmc *MainMsgCache) FlushToFiles() {
	mmc.mu.Lock()
	defer mmc.mu.Unlock()

	for fileID, message := range mmc.cache {
		err := mmc.RetryWrite2File(fileID, message)
		if err != nil {
			log.Printf("Failed to flush messages to file %s: %v", fileID, err)
		} else {
			delete(mmc.cache, fileID)
		}
	}
}

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
			log.Println("Stopping flush to file:", ctx.Err())
			return
		case <-ticker.C:
			log.Println("Ticker ticked, executing task")
			w.task()
		}
	}
}

// AddUsers simulates adding users and messages to message channels
func AddUsers(channels []chan Message, tokens []string, numMsg int) {
	for i := 1; i <= numMsg; i++ {
		for user, ch := range channels {
			msg := Message{
				Token:  tokens[user],
				FileID: fmt.Sprintf("file_%d", user),
				Data:   fmt.Sprintf("Message %d from user %d", i, user),
			}
			ch <- msg
		}
		time.Sleep(1 * time.Second)
	}

	for _, ch := range channels {
		close(ch)
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

	// Handle context cancellation
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := &sync.WaitGroup{}

	// Start worker goroutines
	for i := 0; i < countWorkers; i++ {
		wg.Add(1)
		worker := NewWorker(cache.FlushToFiles, workerInterval)
		go worker.Start(ctx, wg)
	}

	// Start goroutines for message channels
	for _, ch := range messageChannels {
		go func(mc chan Message) {
			for msg := range mc {
				cache.WriteMsgs2Cache(msg)
				log.Println("added to cache:", msg)
			}
		}(ch)
	}

	// Simulate adding users and messages
	go AddUsers(messageChannels, validTokens, numMsg)

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Main: context done, waiting for worker to finish")

	wg.Wait()
	log.Println("Main: worker has been stopped")
}
