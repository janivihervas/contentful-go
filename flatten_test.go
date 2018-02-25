package contentful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToSys(t *testing.T) {
	sys, ok := parseToSys(nil)
	assert.False(t, ok)

	sys, ok = parseToSys(make(map[string]interface{}))
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id": "",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "",
		"linkType": "",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "",
		"linkType": "",
		"type":     "",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "id",
		"linkType": "",
		"type":     "",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "id",
		"linkType": "linkType",
		"type":     "",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "id",
		"linkType": "Asset",
		"type":     "",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "id",
		"linkType": "Asset",
		"type":     "type",
	})
	assert.False(t, ok)

	sys, ok = parseToSys(map[string]interface{}{
		"id":       "id",
		"linkType": "Asset",
		"type":     "Link",
	})
	assert.True(t, ok)
	assert.Equal(t, sys.ID, "id")
	assert.Equal(t, sys.LinkType, "Asset")
	assert.Equal(t, sys.Type, "Link")
}

func TestFetchReferenceWrongIncludeType(t *testing.T) {
	_, err := fetchReference(includes{}, sys{
		LinkType: "linkType",
	})
	assert.NotNil(t, err)
}
