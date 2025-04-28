package main

import (
	// "context"
	"hashing/client"
	"hashing/kubeclient"
	"hashing/metrics"
	"hashing/server"

	// ph "hashing/hashing"
	// "fmt"
	"hashing/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Close()
		os.Exit(0)
	}()
	metrics.Init()
	go server.StartGrpcServer()
	go kubeclient.StartKubeClient()
	var router *mux.Router = router.Router()
	log.Fatal(http.ListenAndServe(":8085", router))
}
