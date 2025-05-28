package main

import (
	"hashing/kubeclient"
	"hashing/router"
	"hashing/server"
	websocketserver "hashing/websocket-server"
	"log"
	"net/http"
)

func main() {
	// Starting grpc server.
	go server.StartGrpcServer()
	// Starting kubernetes watcher.
	go kubeclient.StartKubeClient()
	//
	go websocketserver.GetHubInstance().Run()
	// Starting http server.
	go func() {
		log.Println("http server started on :8085")
		if err := http.ListenAndServe(":8085", router.Router()); err != nil {
			log.Fatalf("http server error: %v", err)
		}
	}()
	// blocking main thread forever.
	select {}
}
