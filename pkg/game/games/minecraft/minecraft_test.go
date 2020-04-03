package minecraft

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

func TestNewGame(t *testing.T) {
	baseConfig := gameinterface.Config{
		InstanceName:  "default",
		GameDirectory: "../../../../test/minecraft",
	}

	g, err := NewGame(baseConfig)
	assert.NoError(t, err)

	_, ok := g.(*Game)
	assert.True(t, ok, "Must be able to assert the Game's type.")
}

// newTestGame returns a Game object for testing purposes.
func newTestGame(rconCommands map[string]string) *Game {
	config := Config{
		rconConstructor: &fakeRconCreator{
			commands: rconCommands,
		},
	}

	g := buildGame(config)
	return g
}

func TestGame_ListPlayers(t *testing.T) {
	cases := []struct {
		listOutput    string
		expectPlayers []string
	}{
		{
			listOutput:    "There are 0/18 players online:",
			expectPlayers: []string{},
		},
		{
			listOutput:    "There are 1/4 players online:SomeUser",
			expectPlayers: []string{"SomeUser"},
		},
		{
			listOutput:    "There are 1/18 players online:SomeUser,anotheruser",
			expectPlayers: []string{"SomeUser", "anotheruser"},
		},
	}

	for _, c := range cases {
		g := newTestGame(
			map[string]string{
				"list": c.listOutput,
			})
		players, err := g.ListPlayers()
		assert.NoError(t, err)
		assert.Equal(t, c.expectPlayers, players)
	}
}
