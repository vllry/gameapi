package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO should this use a mock, to avoid testing individual game initialization here?
func TestNewGame(t *testing.T) {
	cases := []struct {
		name          string
		gameDirectory string
		expectErr     bool
	}{
		{
			name:          "minecraft",
			gameDirectory: "../../test/minecraft",
		},
		{
			name:          "unavailable",
			gameDirectory: "../../test/minecraft",
			expectErr:     true,
		},
	}

	for _, c := range cases {
		game, err := NewGame(c.name, "default", c.gameDirectory)
		if c.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, game)
		}
	}
}
