package main

import (
	"context"
	"goserver/handler"
	metrics "goserver/metric"
	"goserver/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	metrics.Init()
	// Starting background goroutine for pruning old IPs.
	go handler.PruneOldIps()

	// Setting up HTTP server
	// Running HTTP server with graceful termination of 5 second buffer.
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Router(),
	}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Listening for OS signals (ctrl + c).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	// Blocking until signal is received
	<-stop

	// Creating a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")

	// Attempting graceful shutdown.
	// It will serve request after pressing ctrl + c, but till 5 seconds.
	// After it , It will be exited gracefully without any error.
	// But if a request is taking more then 5 seconds then it will be exited forcefully.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", ctx.Err())
	}

	log.Println("Server exited gracefully", ctx.Err())
}
