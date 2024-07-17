package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	requestLimit  = 5
	resetInterval = 1 * time.Minute
	maxClients    = 0
)

var (
	requestCounts = make(map[string]int)
	sem           = make(chan struct{}, maxClients)
	mu            sync.Mutex
)

// rateLimitMiddleware func checks count of request from ip
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		reqCount := requestCounts[ip]
		if reqCount >= requestLimit {
			mu.Unlock()
			http.Error(w, "Request limit exceeded", http.StatusTooManyRequests)
			return
		}
		requestCounts[ip]++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// clientLimitMiddleware func checks how many clients our server has in current time
// Func use Semaphores
func clientLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "Too many clients", http.StatusServiceUnavailable)
		}
	})
}

// resetRequstCounts func starts an infinite loop in goroutine
// then weight and reset the counter every resetInterval
func resetRequstCounts() {
	for {
		time.Sleep(resetInterval)
		mu.Lock()
		requestCounts = make(map[string]int)
		mu.Unlock()
	}
}

// Simple http handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from server!")
}

func main() {
	go resetRequstCounts()

	mux := http.NewServeMux()
	mux.Handle("/", clientLimitMiddleware(rateLimitMiddleware(http.HandlerFunc(helloHandler))))

	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	log.Println("Starting server on :8888")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
