package backup

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"time"

	"github.com/vllry/gameapi/pkg/backup/archive"
	"github.com/vllry/gameapi/pkg/game/identifier"
)

type Manager struct {
	storage Storage
}

// NewManager returns a backup manager using Google Cloud Storage.
func NewManager(storage Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

// Upload uploads the contents of the provided buffer to the storage backend.
func (m *Manager) Upload(ctx context.Context, buf bytes.Buffer, backupTime time.Time, gameIdentifier identifier.GameIdentifier) error {
	fileName := filePathName(gameIdentifier, backupTime)
	return m.storage.Upload(ctx, buf, fileName)
}

// ArchiveDirectory takes a directory path, and returns a buffer containing a targz of that directory.
func ArchiveDirectory(directoryPath string) (bytes.Buffer, error) {
	return archive.ArchiveFilesToBuffer(directoryPath)
}

func filePathName(gameIdentifier identifier.GameIdentifier, t time.Time) string {
	timeString := t.Format(time.RFC1123)
	return path.Join(
		gameIdentifier.Game, gameIdentifier.Instance,
		fmt.Sprintf("%s-%s-%s", gameIdentifier.Game, gameIdentifier.Instance, timeString),
	)
}

func NowUtc() time.Time {
	return time.Now().UTC()
}
