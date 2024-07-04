package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DimTur/multi_user_rw_sys/configs"
	"github.com/DimTur/multi_user_rw_sys/internal/handlers"
	cache "github.com/DimTur/multi_user_rw_sys/internal/service/cache_writer"
	"github.com/DimTur/multi_user_rw_sys/internal/service/validator"
	"github.com/DimTur/multi_user_rw_sys/internal/service/worker"
	"github.com/DimTur/multi_user_rw_sys/models"
)

func main() {
	// Config
	configFile := "../configs/config.yaml"
	config, err := configs.GetConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	workerInterval, err := time.ParseDuration(config.WorkerInterval)
	if err != nil {
		fmt.Printf("Error parsing worker interval: %v\n", err)
		return
	}

	validator := validator.NewMainTokenValidator(config.ValidTokens)
	cache := cache.NewMainMsgCache(validator)

	messageChannels := make([]chan models.Message, config.NumUsers)
	for i := range messageChannels {
		messageChannels[i] = make(chan models.Message, 100)
	}

	// Handle context cancellation
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := &sync.WaitGroup{}

	// Start worker goroutines
	for i := 0; i < config.CountWorkers; i++ {
		wg.Add(1)
		worker := worker.NewWorker(cache.FlushToFiles, workerInterval)
		go worker.Start(ctx, wg)
	}

	// Start goroutines for message channels
	for _, ch := range messageChannels {
		go func(mc chan models.Message) {
			for msg := range mc {
				cache.WriteMsgs2Cache(msg)
				log.Println("added to cache:", msg)
			}
		}(ch)
	}

	// Simulate adding users and messages
	go handlers.AddUsers(messageChannels, config.ValidTokens, config.NumMsg)

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Main: context done, waiting for worker to finish")

	wg.Wait()
	log.Println("Main: worker has been stopped")
}
