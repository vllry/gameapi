package archive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_archive(t *testing.T) {
	err := ArchiveFilesToFile("./test", "./test.tar.gz")
	assert.NoError(t, err)

	_, err = UnArchiveToBytes("./test.tar.gz")
	assert.NoError(t, err)
}
