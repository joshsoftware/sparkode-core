package service

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/joshsoftware/sparkode-core/handler"
)

/* The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/ping", handler.PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/executeCode", handler.RuncodeHandler).Methods(http.MethodPost)

	return
}
