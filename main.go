package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	requestCount uint64
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Do some work here
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	// Set up a health check endpoint
	http.HandleFunc("/healthcheck", handleHealthCheck)

	// Set up mock endpoints that run IRA functions
	http.HandleFunc("/concurrent1", handleConcurrent1)
	http.HandleFunc("/concurrent2", handleConcurrent2)

	// Set up a handler for incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Increment the request counter
		atomic.AddUint64(&requestCount, 1)

		// Set a timeout for the incoming request
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)

		// Set a timeout for the outgoing response
		w.Header().Set("Connection", "close")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// proxy.ServeHTTP(w, r)
	})

	// Set up a server with timeouts
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a separate goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up a signal handler to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Set a timeout for the server to shut down gracefully
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped")
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Perform some health check here
	// If the server is healthy, return a 200 OK status code
	w.WriteHeader(http.StatusOK)
}

func handleConcurrent1(w http.ResponseWriter, r *http.Request) {
	// Start a concurrent function
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Do some work here
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "Concurrent function 1 finished")
	}()
	wg.Wait()
}

func handleConcurrent2(w http.ResponseWriter, r *http.Request) {
	// Start a concurrent function
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Do some work here
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "Concurrent function 1 finished")
	}()
	wg.Wait()
}
