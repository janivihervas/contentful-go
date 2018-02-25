package contentful

import (
	"net/url"
	"strconv"
)

// SearchParameters for GetMany and GetOne functions
type SearchParameters struct {
	url.Values
}

// Parameters returns initialized SearchParameters
func Parameters() SearchParameters {
	return SearchParameters{
		Values: url.Values{},
	}
}

// ByContentType searches by the defined content type
func (p SearchParameters) ByContentType(contentType string) SearchParameters {
	p.Set("content_type", contentType)
	return p
}

// ByFieldValue searches by a field value
func (p SearchParameters) ByFieldValue(fieldName, fieldValue string) SearchParameters {
	p.Add("fields."+fieldName, fieldValue)
	return p
}

// Limit the returned results from Contentful
func (p SearchParameters) Limit(limit int) SearchParameters {
	p.Set("limit", strconv.Itoa(limit))
	return p
}

// Skip n results from Contentful
func (p SearchParameters) Skip(limit int) SearchParameters {
	p.Set("skip", strconv.Itoa(limit))
	return p
}