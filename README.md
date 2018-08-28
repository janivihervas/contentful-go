# contentful-go

[![CircleCI](https://circleci.com/gh/janivihervas/contentful-go.svg?style=svg)](https://circleci.com/gh/janivihervas/contentful-go)
[![codecov](https://codecov.io/gh/janivihervas/contentful-go/branch/master/graph/badge.svg)](https://codecov.io/gh/janivihervas/contentful-go)

[![Go Report Card](https://goreportcard.com/badge/github.com/janivihervas/contentful-go)](https://goreportcard.com/report/github.com/janivihervas/contentful-go)
[![GoDoc](https://godoc.org/github.com/janivihervas/contentful-go?status.svg)](https://godoc.org/github.com/janivihervas/contentful-go)
[![Sponsored](https://img.shields.io/badge/chilicorn-sponsored-brightgreen.svg?logo=data%3Aimage%2Fpng%3Bbase64%2CiVBORw0KGgoAAAANSUhEUgAAAA4AAAAPCAMAAADjyg5GAAABqlBMVEUAAAAzmTM3pEn%2FSTGhVSY4ZD43STdOXk5lSGAyhz41iz8xkz2HUCWFFhTFFRUzZDvbIB00Zzoyfj9zlHY0ZzmMfY0ydT0zjj92l3qjeR3dNSkoZp4ykEAzjT8ylUBlgj0yiT0ymECkwKjWqAyjuqcghpUykD%2BUQCKoQyAHb%2BgylkAyl0EynkEzmkA0mUA3mj86oUg7oUo8n0k%2FS%2Bw%2Fo0xBnE5BpU9Br0ZKo1ZLmFZOjEhesGljuzllqW50tH14aS14qm17mX9%2Bx4GAgUCEx02JySqOvpSXvI%2BYvp2orqmpzeGrQh%2Bsr6yssa2ttK6v0bKxMBy01bm4zLu5yry7yb29x77BzMPCxsLEzMXFxsXGx8fI3PLJ08vKysrKy8rL2s3MzczOH8LR0dHW19bX19fZ2dna2trc3Nzd3d3d3t3f39%2FgtZTg4ODi4uLj4%2BPlGxLl5eXm5ubnRzPn5%2Bfo6Ojp6enqfmzq6urr6%2Bvt7e3t7u3uDwvugwbu7u7v6Obv8fDz8%2FP09PT2igP29vb4%2BPj6y376%2Bu%2F7%2Bfv9%2Ff39%2Fv3%2BkAH%2FAwf%2FtwD%2F9wCyh1KfAAAAKXRSTlMABQ4VGykqLjVCTVNgdXuHj5Kaq62vt77ExNPX2%2Bju8vX6%2Bvr7%2FP7%2B%2FiiUMfUAAADTSURBVAjXBcFRTsIwHAfgX%2FtvOyjdYDUsRkFjTIwkPvjiOTyX9%2FAIJt7BF570BopEdHOOstHS%2BX0s439RGwnfuB5gSFOZAgDqjQOBivtGkCc7j%2B2e8XNzefWSu%2BsZUD1QfoTq0y6mZsUSvIkRoGYnHu6Yc63pDCjiSNE2kYLdCUAWVmK4zsxzO%2BQQFxNs5b479NHXopkbWX9U3PAwWAVSY%2FpZf1udQ7rfUpQ1CzurDPpwo16Ff2cMWjuFHX9qCV0Y0Ok4Jvh63IABUNnktl%2B6sgP%2BARIxSrT%2FMhLlAAAAAElFTkSuQmCC)](http://spiceprogram.org/oss-sponsorship)


Simple Contentful SDK for Go for fetching data. This SDK automatically injects Contentful's entry and asset references to Go structs. See the full documentation in [GoDoc](https://godoc.org/github.com/janivihervas/contentful-go) for more examples.

## Installation

With `go get`:
```
go get github.com/janivihervas/contentful-go
```

## Import

```go
package main

import (
	"github.com/janivihervas/contentful-go/v2"
)

func main() {
	cms := contentful.New(...)
}
```

## Example

```go
// This represents a content model in Contentful
type Page struct {
	// Custom field
	Title  string           `json:"title"`
	// Convenience type for Contentful assets
	Banner contentful.Asset `json:"banner"`
	// This is a reference to many items in Contentful.
	// References will be automatically injected to the struct
	SubPages []Page `json:"subPages"`
	// Inject entry information, which holds e.g. Contentful ID, content type etc
	contentful.Information
}

cms := contentful.New(
	// You can get these values from Space settings -> API keys
	os.Getenv("CONTENTFUL_TOKEN"),
	os.Getenv("CONTENTFUL_SPACE_ID"),
	// Whether to use the preview api or not. Use preview for example on development environment, so you can safely
	// test out the modified or new content before publishing
	true,
)

// Create your result variables
pages := make([]Page, 1)
page := Page{}

// Returns multiple entries.
err := cms.GetMany(
	// Context for cancellation
	context.Background(),
	contentful.SearchParameters{
		// You can use the verbose way of search parameters if you want. See below for convenience functions.
		// Check the docs from https://www.contentful.com/developers/docs/references/content-delivery-api for all the
		// parameters you can use.
		Values: url.Values{
			"content_type": []string{"page"},
		},
	},
	// Pass a reference to the result so the data can be marshaled into it.
	&pages,
)
if err != nil {
	// handle error
}

// Returns exactly one entry
err = cms.GetOne(
	// Context for cancellation
	context.Background(),
	// You can use convenience methods for the parameters
	contentful.
		Parameters().
		ByContentType("page").
		ByFieldValue("title", "Main page"),
	// Pass a reference to the result so the data can be marshaled into it.
	&page,
)
if err != nil {
	// handle error
}

// Multiple results
fmt.Println(len(pages))
fmt.Println(pages[0].Title)
// Additional information
fmt.Println(pages[0].Information.ID)
fmt.Println(pages[0].Information.ContentType)
fmt.Println(pages[0].Information.Revision)
fmt.Println(pages[0].Information.CreatedAt)
fmt.Println(pages[0].Information.UpdatedAt)
fmt.Println(pages[0].Information.Locale)
// Asset has the same information, except for the content type.
// "Information" can be left out
fmt.Println(pages[0].Banner.ID)
fmt.Println(pages[0].Banner.Revision)
fmt.Println(pages[0].Banner.CreatedAt)
fmt.Println(pages[0].Banner.UpdatedAt)
fmt.Println(pages[0].Banner.Locale)
// Asset information
fmt.Println(pages[0].Banner.Title)
fmt.Println(pages[0].Banner.Description)
fmt.Println(pages[0].Banner.File.FileName)
fmt.Println(pages[0].Banner.File.ContentType)
fmt.Println(strings.Split(pages[0].Banner.File.URL, "/")[2]) // Will be in the form of "//images.ctfassets.net/space.id/asset-id/some-id/orange.png"

// One result
fmt.Println(page.Title)
// Additional information
fmt.Println(page.Information.ID)
fmt.Println(page.Information.ContentType)
fmt.Println(page.Information.Revision)
fmt.Println(page.Information.CreatedAt)
fmt.Println(page.Information.UpdatedAt)
fmt.Println(page.Information.Locale)
// Asset has the same information, except for the content type
// "Information" can be left out
fmt.Println(page.Banner.ID)
fmt.Println(page.Banner.Revision)
fmt.Println(page.Banner.CreatedAt)
fmt.Println(page.Banner.UpdatedAt)
fmt.Println(page.Banner.Locale)
// Asset information
fmt.Println(page.Banner.Title)
fmt.Println(page.Banner.Description)
fmt.Println(page.Banner.File.FileName)
fmt.Println(page.Banner.File.ContentType)
fmt.Println(strings.Split(page.Banner.File.URL, "/")[2]) // Will be in the form of "//images.ctfassets.net/space.id/asset-id/some-id/orange.png"
```

## License

[MIT](LICENSE)
