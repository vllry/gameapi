package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

// GameWrapper provides HTTP API responses for game endpoints.
// TODO could this be more generic and automated?
type GameWrapper struct {
	game gameinterface.GenericGame
}

func NewGameWrapper(g gameinterface.GenericGame) GameWrapper {
	return GameWrapper{
		game: g,
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
