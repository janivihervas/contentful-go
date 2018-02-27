// Package contentful provides a Contentful (https://www.contentful.com/) client
package contentful

import "time"

const (
	previewURL = "https://preview.contentful.com"
	cdnURL     = "https://cdn.contentful.com"
)

// Information about the entry or asset
type Information struct {
	ID          string    `json:"contentfulId"`
	ContentType string    `json:"contentfulContentType"`
	Revision    int       `json:"contentfulRevision"`
	CreatedAt   time.Time `json:"contentfulCreatedAt"`
	UpdatedAt   time.Time `json:"contentfulUpdatedAt"`
	Locale      string    `json:"contentfulLocale"`
}

// Asset from Contentful
type Asset struct {
	Information
	Title       string `json:"title"`
	Description string `json:"description"`
	File        File   `json:"file"`
}

// File of an asset
type File struct {
	URL         string `json:"url"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}

// Contentful client for fetching data from Contentful
type Contentful struct {
	token   string
	spaceID string
	url     string
}

// New creates a new Contentful client
func New(token string, spaceID string, preview bool) *Contentful {
	u := cdnURL
	if preview {
		u = previewURL
	}

	return &Contentful{
		token:   token,
		spaceID: spaceID,
		url:     u,
	}
}
