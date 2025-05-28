package router

import (
	"hashing/handler"
	websocketserver "hashing/websocket-server"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	var r *mux.Router = mux.NewRouter()
	r.HandleFunc("/", handler.GetServer)
	r.HandleFunc("/ws", websocketserver.HandleWS)
	return r
}
