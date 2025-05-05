package router

import (
	"hashing/handler"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.GetServer)
	r.HandleFunc("/get", handler.Get)
	r.Path("/metrics").Handler(promhttp.Handler())

	return r
}
