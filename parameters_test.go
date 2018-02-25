package contentful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	params := Parameters()
	assert.NotNil(t, params.Values)
	assert.Equal(t, "", params.Encode())

	params.
		ByContentType("page").
		ByFieldValue("title", "Main page").
		ByFieldValue("description", "Main page description").
		Skip(1).
		Limit(2)

	assert.Equal(t,
		"content_type=page&fields.description=Main+page+description&fields.title=Main+page&limit=2&skip=1",
		params.Encode(),
	)
}
