package game

import (
	"encoding/json"
	"github.com/vllry/gameapi/pkg/game/gameinterface"
	"net/http"
)

// GameWrapper provides HTTP API responses for game endpoints.
type GameWrapper struct {
	game gameinterface.GenericGame
}

func (g *GameWrapper) ListPlayers(w http.ResponseWriter, r *http.Request) {
	result, err := g.game.ListPlayers()
	if err != nil {
		w.WriteHeader(500)
	}

	b, err := json.Marshal(result)
	w.Write(b)
}
