package main

import (
	// "context"
	"hashing/client"
	"hashing/kubeclient"
	"hashing/metrics"
	"hashing/server"
	websocketserver "hashing/websocket-server"

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
	go websocketserver.GetHub().Run()
	var router *mux.Router = router.Router()

	router.HandleFunc("/ws", websocketserver.HandleWS)
	go func() {
		log.Println("WebSocket server started on :8085")
		if err := http.ListenAndServe(":8085", router); err != nil {
			log.Fatalf("WebSocket server error: %v", err)
		}
	}()
	select {}
}
