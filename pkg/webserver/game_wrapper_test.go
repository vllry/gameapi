package webserver

import (
	"reflect"
	"testing"

	"github.com/vllry/gameapi/pkg/game/games/minecraft"
)

// TestGameWrapper_Functions tests, against a sample Game, that all public functions on the wrapper and the game match.
func TestGameWrapper_Functions(t *testing.T) {
	wrapperFuncNames := map[string]struct{}{}
	gameWrapperType := reflect.TypeOf(&GameWrapper{}) // Is this necessary?
	for i := 0; i < gameWrapperType.NumMethod(); i++ {
		method := gameWrapperType.Method(i)
		wrapperFuncNames[method.Name] = struct{}{}
	}

	gameFuncNames := map[string]struct{}{}
	gameType := reflect.TypeOf(&minecraft.Game{})
	for i := 0; i < gameType.NumMethod(); i++ {
		method := gameType.Method(i)
		gameFuncNames[method.Name] = struct{}{}
	}

	if !reflect.DeepEqual(wrapperFuncNames, gameFuncNames) {
		t.Errorf("game object and wrapper object functions didn't match")
	}
}
