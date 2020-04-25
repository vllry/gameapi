package backup

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vllry/gameapi/pkg/game/identifier"
)

// TestManager_Upload is a minimal test of the Upload function, to a fake backend.
func TestManager_Upload(t *testing.T) {
	id := identifier.GameIdentifier{
		Game:     "q",
		Instance: "q",
	}

	m := Manager{storage: &FakeStorage{}}
	err := m.Upload(context.Background(), bytes.Buffer{}, time.Now(), id)

	assert.NoError(t, err)
}
