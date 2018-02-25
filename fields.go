package contentful

import "time"

// ID of the entry
type ID struct {
	ID string `json:"contentful_id"`
}

// ContentType of the entry
type ContentType struct {
	ContentType string `json:"contentful_contentType"`
}

// Revision of the entry
type Revision struct {
	Revision int `json:"contentful_revision"`
}

// CreatedAt of the entry
type CreatedAt struct {
	CreatedAt time.Time `json:"contentful_createdAt"`
}

// UpdatedAt of the entry
type UpdatedAt struct {
	UpdatedAt time.Time `json:"contentful_updatedAt"`
}

// Locale of the entry
type Locale struct {
	Locale string `json:"contentful_locale"`
}
