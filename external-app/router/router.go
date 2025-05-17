package router

import (
	"hashing/handler"
	websocketserver "hashing/websocket-server"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router() *mux.Router {
	var r *mux.Router = mux.NewRouter()
	r.HandleFunc("/", handler.GetServer)
	r.HandleFunc("/ws", websocketserver.HandleWS)
	r.Path("/metrics").Handler(promhttp.Handler())
	return r
}
