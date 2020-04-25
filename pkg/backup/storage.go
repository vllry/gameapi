package backup

import (
	"bytes"
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Storage is an object store provider interface.
type Storage interface {
	Upload(context.Context, bytes.Buffer, string) error
}

type GoogleCloudStorage struct {
	// credentialFilePath is the path to Google Cloud credentials,
	// or empty if the default SDK credential path should be used.
	// https://github.com/googleapis/google-cloud-go#authorization
	credentialFilePath string
}

func NewGoogleCloudStorage(credentialPath string) *GoogleCloudStorage {
	return &GoogleCloudStorage{
		credentialFilePath: credentialPath,
	}
}

// Upload uploads the contents of the provided buffer to a named Google Cloud Storage file.
func (gc *GoogleCloudStorage) Upload(ctx context.Context, buf bytes.Buffer, fileName string) error {
	var client *storage.Client // TODO consider reusing this?
	var clientCreateErr error
	if len(gc.credentialFilePath) == 0 {
		client, clientCreateErr = storage.NewClient(ctx)
	} else {
		client, clientCreateErr = storage.NewClient(ctx, option.WithCredentialsFile(gc.credentialFilePath))
	}

	if clientCreateErr != nil {
		return clientCreateErr
	}

	return gc.uploadToStorage(ctx, client, fileName, buf)
}

// uploadToStorage uploads the contents of the buffer.
// https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go
func (gc *GoogleCloudStorage) uploadToStorage(ctx context.Context, client *storage.Client, objectName string, buffer bytes.Buffer) error {
	bucket := client.Bucket("nebtown-game-backups")

	uploadCtx, cancel := context.WithTimeout(ctx, time.Minute*30)
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
