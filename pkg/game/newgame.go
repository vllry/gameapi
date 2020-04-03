package game

import (
	"errors"
	"fmt"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
	"github.com/vllry/gameapi/pkg/game/games/minecraft"
)

// NewGame creates a new object that satisfies the GenericGame interface.
// It returns an error if the game name is not supported.
func NewGame(gameName string, instanceName string, gameDirectory string) (gameinterface.GenericGame, error) {
	config := gameinterface.Config{
		InstanceName:  instanceName,
		GameDirectory: gameDirectory,
	}

	// Define the signature of a game constructor.
	var constructor func(config gameinterface.Config) (gameinterface.GenericGame, error)

	if gameName == "minecraft" {
		constructor = minecraft.NewGame
	}

	if constructor == nil {
		return nil, errors.New(fmt.Sprintf("game '%s' not found", gameName))
	} else {
		// Call the constructor.
		game, err := constructor(config)
		return game, err
	}
}
