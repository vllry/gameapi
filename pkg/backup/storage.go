package backup

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const googleCloudBucket = "nebtown-game-backups"

// Storage is an object store interface.
type Storage interface {
	upload(context.Context, io.Reader, string) error
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

// upload uploads the contents of the provided buffer to a named Google Cloud Storage file.
// https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go
func (gc *GoogleCloudStorage) upload(ctx context.Context, reader io.Reader, objectName string) error {
	var client *storage.Client // TODO consider reusing this?
	var clientCreateErr error
	if len(gc.credentialFilePath) == 0 {
		client, clientCreateErr = storage.NewClient(ctx)
	} else {
		client, clientCreateErr = storage.NewClient(ctx, option.WithCredentialsFile(gc.credentialFilePath))
	}

	if clientCreateErr != nil {
		return errors.Wrap(clientCreateErr, "couldn't create storage client")
	}

	bucket := client.Bucket(googleCloudBucket)
	uploadCtx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()
	wc := bucket.Object(objectName).NewWriter(uploadCtx)
	if _, err := io.Copy(wc, reader); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
