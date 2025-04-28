package main

import (
	"goserver/handler"
	metrics "goserver/metric"
	"goserver/router"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	metrics.Init()
	go handler.PruneOldIps()
	var router *mux.Router = router.Router()
	log.Fatal(http.ListenAndServe(":8080", router))
}
