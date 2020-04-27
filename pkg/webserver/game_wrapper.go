package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

// GameWrapper provides HTTP API responses for game endpoints.
type GameWrapper struct {
	game gameinterface.GenericGame
}

func NewGameWrapper(g gameinterface.GenericGame) GameWrapper {
	return GameWrapper{
		game: g,
	}
}

func (g *GameWrapper) Backup(w http.ResponseWriter, r *http.Request) {
	err := g.game.Backup()
	if err != nil {
		log.Println("Backup failed: ", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("done"))
	}
}

func (g *GameWrapper) GetLogs(w http.ResponseWriter, r *http.Request) {
	result, err := g.game.GetLogs()
	if err != nil {
		w.WriteHeader(500)
	}

	b, err := json.Marshal(result)
	w.Write(b)
}

func (g *GameWrapper) ListPlayers(w http.ResponseWriter, r *http.Request) {
	result, err := g.game.ListPlayers()
	if err != nil {
		w.WriteHeader(500)
	}

	b, err := json.Marshal(result)
	w.Write(b)
}
