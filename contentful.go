// Package contentful provides a Contentful (https://www.contentful.com/) client
package contentful

// Contentful client for fetching data from Contentful
type Contentful interface {
}

type contentful struct {
}

// New creates a new Contentful client
func New() Contentful {
	return &contentful{}
}
