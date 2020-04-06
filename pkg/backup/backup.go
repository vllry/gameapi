package backup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/vllry/gameapi/pkg/game/identifier"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/vllry/gameapi/pkg/backup/archive"
)

type Manager struct {
	// credentialFilePath is the path to Google Cloud credentials,
	// or empty if the default SDK credential path should be used.
	// https://github.com/googleapis/google-cloud-go#authorization
	credentialFilePath string
}

func NewManager(credentialFilePath string) *Manager {
	return &Manager{
		credentialFilePath: credentialFilePath,
	}
}

// Backup backs up a directory, and uploads an archive version.
func (m *Manager) Backup(ctx context.Context, directoryPath string, gameName string, gameInstance string) error {
	var client *storage.Client // TODO consider reusing this?
	var clientCreateErr error
	if len(m.credentialFilePath) == 0 {
		client, clientCreateErr = storage.NewClient(ctx)
	} else {
		client, clientCreateErr = storage.NewClient(ctx, option.WithCredentialsFile(m.credentialFilePath))
	}

	if clientCreateErr != nil {
		return clientCreateErr
	}

	buf, err := ArchiveDirectory(directoryPath)
	if err != nil {
		return err
	}

	timeString := time.Now().UTC().Format(time.RFC1123)
	storagePath := path.Join(gameName, gameInstance, fmt.Sprintf("%s-%s-%s", gameName, gameInstance, timeString))
	return uploadToStorage(ctx, client, storagePath, buf)
}

// ArchiveDirectory takes a directory path, and returns a buffer containing a targz of that directory.
func ArchiveDirectory(directoryPath string) (bytes.Buffer, error) {
	return archive.ArchiveFilesToBuffer(directoryPath)
}

// Upload uploads the contents of the provided buffer to a named Google Cloud Storage file.
func (m *Manager) Upload(ctx context.Context, buf bytes.Buffer, backupTime time.Time, gameIdentifier identifier.GameIdentifier) error {
	var client *storage.Client // TODO consider reusing this?
	var clientCreateErr error
	if len(m.credentialFilePath) == 0 {
		client, clientCreateErr = storage.NewClient(ctx)
	} else {
		client, clientCreateErr = storage.NewClient(ctx, option.WithCredentialsFile(m.credentialFilePath))
	}

	if clientCreateErr != nil {
		return clientCreateErr
	}

	storagePath := filePathName(gameIdentifier, backupTime)
	return uploadToStorage(ctx, client, storagePath, buf)
}

// uploadToStorage uploads the contents of the buffer.
// https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go
func uploadToStorage(ctx context.Context, client *storage.Client, objectName string, buffer bytes.Buffer) error {
	bucket := client.Bucket("nebtown-game-backups")

	uploadCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	wc := bucket.Object(objectName).NewWriter(uploadCtx)
	if _, err := io.Copy(wc, &buffer); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
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
