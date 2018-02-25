package contentful_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"net/url"

	"github.com/janivihervas/contentful-go"
)

func Example() {
	// This represents a content model in Contentful
	type Page struct {
		// Custom field
		Title string `json:"title"`
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
	fmt.Println(strings.Split(pages[0].Banner.File.URL, "/")[2]) // Will be in the form of "//images.contentful.com/space.id/asset-id/some-id/orange.png"

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
	fmt.Println(strings.Split(page.Banner.File.URL, "/")[2]) // Will be in the form of "//images.contentful.com/space.id/asset-id/some-id/orange.png"
	// Output:
	// 3
	// Not published page
	// 5CVt4s6uvS0cuym4wmWg2k
	// page
	// 1
	// 2018-02-20 18:19:01.359 +0000 UTC
	// 2018-02-20 18:23:36.818 +0000 UTC
	// en-US
	// 3hBPd2try88AaWsMMU2Ga
	// 1
	// 2018-02-20 18:17:38.568 +0000 UTC
	// 2018-02-20 18:17:58.035 +0000 UTC
	// en-US
	// Blue
	// Blue image
	// blue.png
	// image/png
	// images.contentful.com
	// Main page
	// 2Cbt07njicqO4wSYCQ8CeK
	// page
	// 2
	// 2018-02-20 18:14:49.006 +0000 UTC
	// 2018-02-20 18:24:07.281 +0000 UTC
	// en-US
	// 2BNT5Xj0CsgUOSMkKysYKq
	// 2
	// 2018-02-20 18:14:18.301 +0000 UTC
	// 2018-02-20 18:17:30.57 +0000 UTC
	// en-US
	// Green
	// Green image
	// green.png
	// image/png
	// images.contentful.com
}

func ExampleContentful_GetMany() {
	type Page struct {
		Title    string           `json:"title"`
		Banner   contentful.Asset `json:"banner"`
		SubPages []Page           `json:"subPages"`
	}

	cms := contentful.New(
		os.Getenv("CONTENTFUL_TOKEN"),
		os.Getenv("CONTENTFUL_SPACE_ID"),
		true,
	)
	ctx := context.Background()

	pages := make([]Page, 1)
	params := contentful.Parameters().ByContentType("page")

	err := cms.GetMany(ctx, params, &pages)
	if err != nil {
		// handle error
	}

	fmt.Println(len(pages))
	fmt.Println(pages[0].Title)
	fmt.Println(pages[0].Banner.Title)
	fmt.Println(pages[0].Banner.Description)
	fmt.Println(strings.Split(pages[0].Banner.File.URL, "/")[2])
	// Output:
	// 3
	// Not published page
	// Blue
	// Blue image
	// images.contentful.com
}

func ExampleContentful_GetOne() {
	type Page struct {
		Title    string           `json:"title"`
		Banner   contentful.Asset `json:"banner"`
		SubPages []Page           `json:"subPages"`
	}

	cms := contentful.New(
		os.Getenv("CONTENTFUL_TOKEN"),
		os.Getenv("CONTENTFUL_SPACE_ID"),
		true,
	)
	ctx := context.Background()

	page := Page{}
	params := contentful.Parameters().ByContentType("page").ByFieldValue("title", "Main page")

	err := cms.GetOne(ctx, params, &page)
	if err != nil {
		// handle error
	}

	fmt.Println(page.Title)
	fmt.Println(page.Banner.Title)
	fmt.Println(page.Banner.Description)
	fmt.Println(strings.Split(page.Banner.File.URL, "/")[2])
	fmt.Println(len(page.SubPages))
	// Output:
	// Main page
	// Green
	// Green image
	// images.contentful.com
	// 2
}

func ExampleParameters() {
	type Page struct {
		Title string `json:"title"`
	}

	cms := contentful.New(
		os.Getenv("CONTENTFUL_TOKEN"),
		os.Getenv("CONTENTFUL_SPACE_ID"),
		true,
	)
	ctx := context.Background()

	page := Page{}
	params := contentful.Parameters().
		ByFieldValue("title", "Sub page")

	err := cms.GetOne(ctx, params, &page)
	if err != nil {
		fmt.Println("ByFieldValue requires that ByContentType is called:", err)
	}

	params = params.ByContentType("page")

	err = cms.GetOne(ctx, params, &page)
	if err != nil {
		// Handle error
	}

	fmt.Println(page.Title)
	// Output:
	// ByFieldValue requires that ByContentType is called: non-ok status code: 400
	// Sub page
}

func ExampleInformation() {
	type Page struct {
		Title string `json:"title"`
		// Helper types you can use to access the entry's information
		contentful.Information
	}

	cms := contentful.New(
		os.Getenv("CONTENTFUL_TOKEN"),
		os.Getenv("CONTENTFUL_SPACE_ID"),
		true,
	)
	ctx := context.Background()

	page := Page{}
	params := contentful.Parameters().
		ByContentType("page").
		ByFieldValue("title", "Sub page")

	err := cms.GetOne(ctx, params, &page)
	if err != nil {
		// Handle error
	}

	fmt.Println(page.ID)          // Or page.Information.ID
	fmt.Println(page.ContentType) // Or page.Information.ContentType
	fmt.Println(page.Revision)    // Or page.Information.Revision
	fmt.Println(page.CreatedAt)   // Or page.Information.CreatedAt
	fmt.Println(page.UpdatedAt)   // Or page.Information.UpdatedAt
	fmt.Println(page.Locale)      // Or page.Information.Locale
	// Output:
	// FcAxxzogmsOMcc0kac6Iu
	// page
	// 1
	// 2018-02-20 18:15:09.146 +0000 UTC
	// 2018-02-20 18:19:32.99 +0000 UTC
	// en-US
}

func ExampleAsset() {
	type Page struct {
		Banner contentful.Asset `json:"banner"`
	}

	cms := contentful.New(
		os.Getenv("CONTENTFUL_TOKEN"),
		os.Getenv("CONTENTFUL_SPACE_ID"),
		true,
	)
	ctx := context.Background()

	page := Page{}
	params := contentful.Parameters().
		ByContentType("page").
		ByFieldValue("title", "Main page")

	err := cms.GetOne(ctx, params, &page)
	if err != nil {
		// Handle error
	}

	fmt.Println(page.Banner.Title)
	fmt.Println(page.Banner.Description)
	// Asset has also access to additional information from Contentful
	fmt.Println(page.Banner.ID)          // Or page.Banner.Information.ID
	fmt.Println(page.Banner.ContentType) // Or page.Banner.Information.ContentType // Assets don't have a content type
	fmt.Println(page.Banner.Revision)    // Or page.Banner.Information.Revision
	fmt.Println(page.Banner.CreatedAt)   // Or page.Banner.Information.CreatedAt
	fmt.Println(page.Banner.UpdatedAt)   // Or page.Banner.Information.UpdatedAt
	fmt.Println(page.Banner.Locale)      // Or page.Banner.Information.Locale
	// Output:
	// Green
	// Green image
	// 2BNT5Xj0CsgUOSMkKysYKq
	//
	// 2
	// 2018-02-20 18:14:18.301 +0000 UTC
	// 2018-02-20 18:17:30.57 +0000 UTC
	// en-US
}

func ExampleFile() {
	type Page struct {
		Banner contentful.Asset `json:"banner"`
	}

	cms := contentful.New(
		os.Getenv("CONTENTFUL_TOKEN"),
		os.Getenv("CONTENTFUL_SPACE_ID"),
		true,
	)
	ctx := context.Background()

	page := Page{}
	params := contentful.Parameters().
		ByContentType("page").
		ByFieldValue("title", "Main page")

	err := cms.GetOne(ctx, params, &page)
	if err != nil {
		// Handle error
	}

	fmt.Println(page.Banner.File.FileName)
	fmt.Println(page.Banner.File.ContentType)
	fmt.Println(strings.Split(page.Banner.File.URL, "/")[2])
	// Output:
	// green.png
	// image/png
	// images.contentful.com
}
