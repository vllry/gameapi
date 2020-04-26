package backup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vllry/gameapi/pkg/game/identifier"
)

// TestManager_Upload is a minimal test of the upload function, to a fake backend.
func TestBackup_Upload(t *testing.T) {
	b := Backup{
		gameId: identifier.GameIdentifier{
			Game:     "q",
			Instance: "q",
		},
		storage:         &FakeStorage{},
		archiveFilePath: "storage_test.go", // Throwing a random file at this function is not representative of expected use.
	}
	err := b.Upload(context.Background())

	assert.NoError(t, err)
}
