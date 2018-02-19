package contentful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	cms := New("token", "space", false)

	c, ok := cms.(*client)

	assert.True(ok)
	assert.Equal("token", c.token)
	assert.Equal("space", c.spaceID)
	assert.Equal(cdnURL, c.url)

	cms = New("token", "space", true)

	c, ok = cms.(*client)

	assert.True(ok)
	assert.Equal(previewURL, c.url)
}
