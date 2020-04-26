package backup

import (
	"context"
	"io"
)

type FakeStorage struct{}

// upload uploads the contents of the provided buffer to a named Google Cloud Storage file.
func (f *FakeStorage) upload(ctx context.Context, reader io.Reader, fileName string) error {
	return nil
}
