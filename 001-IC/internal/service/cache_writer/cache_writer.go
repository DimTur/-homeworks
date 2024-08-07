package cachewriter

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/DimTur/multi_user_rw_sys/internal/service/validator"
	"github.com/DimTur/multi_user_rw_sys/models"
)

// MessageCache interface for message caching
//
//go:generate mockgen -source=cache_writer.go -destination=cachewriter_mock.go -package cachewriter
type MessageCache interface {
	WriteMsgs2Cache(msg models.Message)
	FlushToFiles()
}

// MainMsgCache caches messages and flushes them to files
type MainMsgCache struct {
	Mu        sync.RWMutex
	Cache     map[string][]models.Message
	Validator validator.TokenValidator
}

// NewMainMsgCache creates a new instance of MainMsgCache
func NewMainMsgCache(validator validator.TokenValidator) *MainMsgCache {
	return &MainMsgCache{
		Cache:     make(map[string][]models.Message),
		Validator: validator,
	}
}

// WriteMsgs2Cache adds a message to the cache
func (mmc *MainMsgCache) WriteMsgs2Cache(msg models.Message) {
	mmc.Mu.Lock()
	defer mmc.Mu.Unlock()

	if mmc.Validator.Validate(msg.Token) {
		mmc.Cache[msg.FileID] = append(mmc.Cache[msg.FileID], msg)
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing file: ", err)
	}
}

// RetryWrite2File retries writing messages to a file with exponential backoff in case of failures.
// It attempts multiple times and returns an error if unsuccessful after the maximum retries.
func (mmc *MainMsgCache) RetryWrite2File(fileID string, messages []models.Message) error {
	const maxRetries = 5
	const baseDelay = 3 * time.Second

	filePath := fmt.Sprintf("../data/%s.txt", fileID)
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
	mmc.Mu.Lock()
	defer mmc.Mu.Unlock()

	for fileID, message := range mmc.Cache {
		err := mmc.RetryWrite2File(fileID, message)
		if err != nil {
			log.Printf("Failed to flush messages to file %s: %v", fileID, err)
		} else {
			delete(mmc.Cache, fileID)
		}
	}
}
