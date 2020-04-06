package backup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vllry/gameapi/pkg/game/identifier"
)

func Test_filePathName(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC1123, "Mon, 02 Jan 2006 15:04:05 UTC")
	assert.NoError(t, err)

	actual := filePathName(
		identifier.GameIdentifier{
			"gmod2",
			"beetown",
		},
		parsedTime,
	)

	assert.Equal(t, "gmod2/beetown/gmod2-beetown-Mon, 02 Jan 2006 15:04:05 UTC", actual)
}
