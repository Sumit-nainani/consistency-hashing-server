package router

import (
	"goserver/handler"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.GetServer)
	r.Path("/metrics").Handler(promhttp.Handler())
	return r
}
