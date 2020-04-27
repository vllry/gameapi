package webserver

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

func Start(g gameinterface.GenericGame) {
	gameWrapper := NewGameWrapper(g)

	router := mux.NewRouter()

	router.HandleFunc("/", Index)
	router.HandleFunc("/backup", gameWrapper.Backup)
	router.HandleFunc("/players", gameWrapper.ListPlayers)
	router.HandleFunc("/logs", gameWrapper.GetLogs)

	http.Handle("/", router)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting webserver...")
	log.Fatal(srv.ListenAndServe())
}
