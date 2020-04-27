package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/pkg/errors"

	"github.com/vllry/gameapi/pkg/game/identifier"
)

// Manager is used to handle backups, and backup uploading.
// It contains a storage instance.
type Manager struct {
	tempDirectory string
	storage       Storage
}

// Backup represents a single backup operation.
type Backup struct {
	archiveFilePath string
	gameId          identifier.GameIdentifier
	storage         Storage
}

// NewManager returns a backup manager using Google Cloud Storage.
func NewManager(storage Storage, tempDirectory string) *Manager {
	return &Manager{
		tempDirectory: tempDirectory,
		storage:       storage,
	}
}

// NewTestManager returns a manager for use in tests.
// It is exported to allow other modules to write tests involving backups.
func NewTestManager() *Manager {
	return NewManager(&FakeStorage{}, "../../../../_tmp") // Relative path based on game directory structure.
}

// Upload uses the storage provider's upload function to upload the backup file.
func (b *Backup) Upload(ctx context.Context) error {
	f, err := os.Open(b.archiveFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read archive file")
	}
	defer f.Close()

	err = b.storage.upload(
		ctx,
		f,
		path.Join(storagePathPrefix(b.gameId), path.Base(b.archiveFilePath)), // Combine storage prefix and file name.
	)
	return errors.Wrap(err, "failed to upload archive")
}

// ArchiveDirectory creates a 7zip of the specified directory, and returns the path to the 7zip file.
func (m *Manager) ArchiveDirectory(directoryPrefix string, fileNames []string, gameId identifier.GameIdentifier) (Backup, error) {
	outputFile := path.Join(m.tempDirectory, archiveFileName(gameId, time.Now().UTC()))

	for i := range fileNames {
		fileNames[i] = path.Join(directoryPrefix, fileNames[i])
	}

	args := []string{"a", outputFile}
	args = append(args, fileNames...)
	log.Println("Running: 7zr", args)
	cmd := exec.Command("7zr", args...)
	err := cmd.Run()
	return Backup{
		archiveFilePath: outputFile,
		gameId:          gameId,
		storage:         m.storage,
	}, errors.Wrap(err, "failed to execute archive command")
}

// archiveFileName returns a file name for the archive.
func archiveFileName(gameIdentifier identifier.GameIdentifier, t time.Time) string {
	timeString := t.Format(time.RFC1123)
	return fmt.Sprintf("%s-%s-%s.7z", gameIdentifier.Game, gameIdentifier.Instance, timeString)
}

// storagePathPrefix returns a path prefix to store an archive file.
func storagePathPrefix(gameId identifier.GameIdentifier) string {
	return path.Join(gameId.Game, gameId.Instance)
}
