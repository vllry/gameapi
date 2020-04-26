package minecraft

import (
	"testing"

	"github.com/vllry/gameapi/pkg/backup"

	"github.com/vllry/gameapi/pkg/game/identifier"

	"github.com/stretchr/testify/assert"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

func TestNewGame(t *testing.T) {
	baseConfig := gameinterface.Config{
		Identifier: identifier.GameIdentifier{
			Game:     "minecraft",
			Instance: "default",
		},
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
		base: gameinterface.Config{
			Identifier: identifier.GameIdentifier{
				Game:     "minecraft",
				Instance: "default",
			},
			GameDirectory: "../../../../test/minecraft",
			BackupManager: backup.NewTestManager(),
		},
		rconConstructor: &fakeRconCreator{
			commands: rconCommands,
		},
	}

	g := buildGame(config)
	return g
}

func TestGame_Backup(t *testing.T) {
	g := newTestGame(map[string]string{
		"save-off": "",
		"save":     "",
		"save-on":  "",
	})

	err := g.Backup()
	assert.NoError(t, err)
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

func TestGame_GetLogs(t *testing.T) {
	expectLogs := `This is log line 1
This thing ain't on autopilot.
End.`

	g := newTestGame(map[string]string{})
	logs, err := g.GetLogs()
	assert.NoError(t, err)
	assert.Equal(t, expectLogs, logs)
}
