package contentful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	cms := New("token", "space", false)

	assert.Equal("token", cms.token)
	assert.Equal("space", cms.spaceID)
	assert.Equal(cdnURL, cms.url)

	cms = New("token", "space", true)
	assert.Equal(previewURL, cms.url)
}
