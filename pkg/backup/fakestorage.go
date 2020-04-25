package backup

import (
	"bytes"
	"context"
)

type FakeStorage struct{}

// Upload uploads the contents of the provided buffer to a named Google Cloud Storage file.
func (f *FakeStorage) Upload(ctx context.Context, buf bytes.Buffer, fileName string) error {
	return nil
}
