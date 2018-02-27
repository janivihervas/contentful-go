package contentful

import (
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestContentful_GetServerFails(t *testing.T) {
	t.Parallel()
	var (
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
		}
		ctx    = context.Background()
		result = make([]map[string]interface{}, 1)
	)

	err := cms.GetMany(ctx, Parameters(), &result)
	assert.NotNil(t, err)

	err = cms.GetOne(ctx, Parameters(), &result)
	assert.NotNil(t, err)
}

func TestContentful_GetWrongTotal(t *testing.T) {
	t.Parallel()

	var (
		response = searchResults{
			Total: 0,
		}
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(response)
			assert.Nil(t, err)
		}))
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
			url:     server.URL,
		}
		ctx    = context.Background()
		result = make([]map[string]interface{}, 1)
	)
	defer server.Close()

	err := cms.GetMany(ctx, Parameters(), &result)
	assert.NotNil(t, err)
	err = cms.GetOne(ctx, Parameters(), &result)
	assert.NotNil(t, err)

	response.Total = 2
	err = cms.GetOne(ctx, Parameters(), &result)
	assert.NotNil(t, err)

	response.Total = 1
	err = cms.GetMany(ctx, Parameters(), &result)
	assert.NotNil(t, err)
	err = cms.GetOne(ctx, Parameters(), &result)
	assert.NotNil(t, err)
}

func TestContentful_GetParseFails(t *testing.T) {
	t.Parallel()

	bytes, err := ioutil.ReadFile("test_data/reference_fields.json")
	assert.Nil(t, err)
	fields := make(map[string]interface{})
	err = json.Unmarshal(bytes, &fields)
	assert.Nil(t, err)

	var (
		response = searchResults{
			Total: 1,
			Items: []item{
				{
					Fields: fields,
				},
			},
		}
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(response)
			assert.Nil(t, err)
		}))
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
			url:     server.URL,
		}
		ctx        = context.Background()
		resultMany = make([]map[string]interface{}, 1)
		resultOne  = make(map[string]interface{})
	)
	defer server.Close()

	err = cms.GetMany(ctx, Parameters(), &resultMany)
	assert.NotNil(t, err)
	err = cms.GetOne(ctx, Parameters(), &resultOne)
	assert.NotNil(t, err)
}

func TestContentful_GetUnMarshalFails(t *testing.T) {
	t.Parallel()

	var (
		dataFile = ""
		server   = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("test_data/" + dataFile)
			assert.Nil(t, err)
			_, err = w.Write(bytes)
			assert.Nil(t, err)
		}))
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
			url:     server.URL,
		}
		ctx        = context.Background()
		resultMany = make([]map[string]interface{}, 1)
		resultOne  = make(map[string]interface{})
	)
	defer server.Close()

	dataFile = "preview_all_pages.json"
	err := cms.GetMany(ctx, Parameters(), &resultOne)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "json: cannot unmarshal")

	dataFile = "preview_main_page.json"
	err = cms.GetOne(ctx, Parameters(), &resultMany)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "json: cannot unmarshal")
}

func TestContentful_Get(t *testing.T) {
	t.Parallel()

	var (
		dataFile = ""
		server   = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("test_data/" + dataFile)
			assert.Nil(t, err)
			_, err = w.Write(bytes)
			assert.Nil(t, err)
		}))
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
			url:     server.URL,
		}
		ctx = context.Background()
	)
	defer server.Close()

	resultMany := make([]map[string]interface{}, 1)
	dataFile = "preview_all_pages.json"
	err := cms.GetMany(ctx, Parameters(), &resultMany)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(resultMany))

	cases := []struct {
		expected interface{}
		actual   interface{}
	}{
		{3, len(resultMany)},
		{"Sub page", resultMany[0]["title"]},
		{"page", resultMany[0]["contentfulContentType"]},
		{"FcAxxzogmsOMcc0kac6Iu", resultMany[0]["contentfulId"]},
		{"en-US", resultMany[0]["contentfulLocale"]},
		{"Not published page", resultMany[1]["title"]},
		{"page", resultMany[1]["contentfulContentType"]},
		{"5CVt4s6uvS0cuym4wmWg2k", resultMany[1]["contentfulId"]},
		{"en-US", resultMany[1]["contentfulLocale"]},
		{"Main page", resultMany[2]["title"]},
		{"page", resultMany[2]["contentfulContentType"]},
		{"2Cbt07njicqO4wSYCQ8CeK", resultMany[2]["contentfulId"]},
		{"en-US", resultMany[2]["contentfulLocale"]},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual)
	}

	resultOne := make(map[string]interface{})
	dataFile = "preview_all_pages.json"
	err = cms.GetOne(ctx, Parameters(), &resultOne)
	assert.NotNil(t, err)

	dataFile = "preview_main_page.json"
	err = cms.GetOne(ctx, Parameters(), &resultOne)
	assert.Nil(t, err)

	cases = []struct {
		expected interface{}
		actual   interface{}
	}{
		{2, len(resultOne["subPages"].([]interface{}))},
		{"Main page", resultOne["title"]},
		{"page", resultOne["contentfulContentType"]},
		{"2Cbt07njicqO4wSYCQ8CeK", resultOne["contentfulId"]},
		{"en-US", resultOne["contentfulLocale"]},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.actual)
	}
	assert.Equal(t, "Main page", resultOne["title"])
	assert.Equal(t, 2, len(resultOne["subPages"].([]interface{})))
}

func TestContentful_search(t *testing.T) {
	t.Parallel()

	var (
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
		}
		ctx = context.Background()
	)

	t.Run("Should return an error if no url is set", func(tt *testing.T) {
		_, err := cms.search(ctx, SearchParameters{})
		assert.NotNil(tt, err)
	})

	t.Run("Client should call the correct endpoint with correct bearer token", func(tt *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(tt, "/spaces/spaceID/entries?include=10", r.URL.Path+"?"+r.URL.RawQuery)
			assert.Equal(tt, "Bearer token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cms.url = server.URL
		_, _ = cms.search(ctx, Parameters())
	})

	t.Run("Should return an error if Contentful responds with non-200 status code", func(tt *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cms.url = server.URL
		_, err := cms.search(ctx, Parameters())
		assert.NotNil(tt, err)
		assert.Contains(tt, err.Error(), "500")
	})

	t.Run("Should return an error if body can't be read", func(tt *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
		}))
		defer server.Close()

		cms.url = server.URL
		_, err := cms.search(ctx, Parameters())
		assert.NotNil(tt, err)
	})

	t.Run("Should not return an error and should return correctly parsed search results", func(tt *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("test_data/prod_all_pages.json")
			assert.Nil(tt, err)
			_, err = w.Write(bytes)
			assert.Nil(tt, err)
		}))
		defer server.Close()

		cms.url = server.URL
		response, err := cms.search(ctx, Parameters())
		assert.Nil(tt, err)
		assert.Equal(tt, 2, response.Total)
		assert.Equal(tt, 0, response.Skip)
		assert.Equal(tt, 100, response.Limit)
		assert.Equal(tt, 2, len(response.Items))
	})
}

func TestAppendIncludes(t *testing.T) {
	t.Parallel()

	response := searchResults{
		Items: []item{},
		Includes: includes{
			Entry: []item{
				{
					Sys: itemInfo{
						ID: "1",
					},
				},
			},
		},
	}

	appendIncludes(&response)
	assert.Equal(t, 1, len(response.Includes.Entry))

	response.Items = []item{
		{
			Sys: itemInfo{
				Type: linkTypeAsset,
			},
		},
	}
	appendIncludes(&response)
	assert.Equal(t, 1, len(response.Includes.Entry))

	response.Items = append(response.Items, item{
		Sys: itemInfo{
			ID:   "2",
			Type: linkTypeEntry,
		},
	})
	appendIncludes(&response)
	assert.Equal(t, 2, len(response.Includes.Entry))
	assert.Equal(t, "1", response.Includes.Entry[0].Sys.ID)
	assert.Equal(t, "2", response.Includes.Entry[1].Sys.ID)
}
