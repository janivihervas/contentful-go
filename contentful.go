// Package contentful provides a Contentful (https://www.contentful.com/) client
package contentful

import "net/url"

const (
	previewURL = "https://preview.client.com"
	cdnURL     = "https://cdn.client.com"
)

// Contentful client for fetching data from Contentful
type Contentful interface {
	Search(parameters url.Values, data interface{}) error
}

type client struct {
	token   string
	spaceID string
	url     string
}

// New creates a new Contentful client
func New(token string, spaceID string, preview bool) Contentful {
	u := cdnURL
	if preview {
		u = previewURL
	}

	return &client{
		token:   token,
		spaceID: spaceID,
		url:     u,
	}
}
