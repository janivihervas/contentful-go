package contentful_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/janivihervas/contentful-go"
)

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

	fmt.Println("Pages returned:", len(pages))
	fmt.Println("Title:", pages[0].Title)
	fmt.Println("Banner title:", pages[0].Banner.Title)
	fmt.Println("Banner description:", pages[0].Banner.Description)
	fmt.Println("Banner file url host:", strings.Split(pages[0].Banner.File.URL, "/")[2])
	// Output:
	// Pages returned: 3
	// Title: Not published page
	// Banner title: Blue
	// Banner description: Blue image
	// Banner file url host: images.contentful.com
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
