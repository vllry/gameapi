package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func startWebserver(game gameinterface.GenericGame) {
	router := mux.NewRouter()

	router.HandleFunc("/", Index)
	router.HandleFunc("/players", Index)

	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
