package contentful

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	assert.Error(t, err)

	err = cms.GetOne(ctx, Parameters(), &result)
	assert.Error(t, err)
}

func TestContentful_GetWrongTotal(t *testing.T) {
	t.Parallel()

	var (
		items = item{Fields: map[string]interface{}{
			"foo": "bar",
		}}
		response = searchResults{}
		server   = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(response)
			assert.NoError(t, err)
		}))
		cms = Contentful{
			token:   "token",
			spaceID: "spaceID",
			url:     server.URL,
		}
		ctx       = context.Background()
		result    = make([]map[string]interface{}, 1)
		resultOne = make(map[string]interface{})
	)
	defer server.Close()

	t.Run("Zero entries", func(t *testing.T) {
		response.Total = 0
		response.Items = make([]item, 0)
		err := cms.GetMany(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrNoEntries, err)
		err = cms.GetOne(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrNoEntries, err)
	})

	t.Run("Two entries", func(t *testing.T) {
		response.Total = 2
		response.Items = []item{items, items}
		err := cms.GetMany(ctx, Parameters(), &result)
		assert.NoError(t, err)
		err = cms.GetOne(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrMoreThanOneEntry, err)
	})

	t.Run("One entry", func(t *testing.T) {
		response.Total = 1
		response.Items = []item{items}
		err := cms.GetMany(ctx, Parameters(), &result)
		assert.NoError(t, err)
		err = cms.GetOne(ctx, Parameters(), &resultOne)
		assert.NoError(t, err)
	})

	t.Run("One in total, zero in items", func(t *testing.T) {
		response.Total = 1
		response.Items = make([]item, 0)
		err := cms.GetMany(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrNoEntries, err)
		err = cms.GetOne(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrNoEntries, err)
	})

	t.Run("Zero in total, one in items", func(t *testing.T) {
		response.Total = 0
		response.Items = []item{items}
		err := cms.GetMany(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrNoEntries, err)
		err = cms.GetOne(ctx, Parameters(), &result)
		assert.Error(t, err)
		assert.Equal(t, ErrNoEntries, err)
	})
}

func TestContentful_GetParseFails(t *testing.T) {
	t.Parallel()

	bytes, err := ioutil.ReadFile("testdata/reference_fields.json")
	assert.NoError(t, err)
	fields := make(map[string]interface{})
	err = json.Unmarshal(bytes, &fields)
	assert.NoError(t, err)

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
			assert.NoError(t, err)
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
	assert.Error(t, err)
	err = cms.GetOne(ctx, Parameters(), &resultOne)
	assert.Error(t, err)
}

func TestContentful_GetUnMarshalFails(t *testing.T) {
	t.Parallel()

	var (
		dataFile = ""
		server   = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("testdata/" + dataFile)
			assert.NoError(t, err)
			_, err = w.Write(bytes)
			assert.NoError(t, err)
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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: cannot unmarshal")

	dataFile = "preview_main_page.json"
	err = cms.GetOne(ctx, Parameters(), &resultMany)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: cannot unmarshal")
}

func TestContentful_Get(t *testing.T) {
	t.Parallel()

	var (
		dataFile = ""
		server   = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("testdata/" + dataFile)
			assert.NoError(t, err)
			_, err = w.Write(bytes)
			assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.Error(t, err)

	dataFile = "preview_main_page.json"
	err = cms.GetOne(ctx, Parameters(), &resultOne)
	assert.NoError(t, err)

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

	t.Run("Should return an error if no url is set", func(t *testing.T) {
		_, err := cms.search(ctx, SearchParameters{})
		assert.Error(t, err)
	})

	t.Run("Client should call the correct endpoint with correct bearer token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/spaces/spaceID/entries?include=10", r.URL.Path+"?"+r.URL.RawQuery)
			assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cms.url = server.URL
		_, _ = cms.search(ctx, Parameters())
	})

	t.Run("Should return an error if Contentful responds with non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cms.url = server.URL
		_, err := cms.search(ctx, Parameters())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("Should return an error if body can't be read", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
		}))
		defer server.Close()

		cms.url = server.URL
		_, err := cms.search(ctx, Parameters())
		assert.Error(t, err)
	})

	t.Run("Should not return an error and should return correctly parsed search results", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("testdata/prod_all_pages.json")
			assert.NoError(t, err)
			_, err = w.Write(bytes)
			assert.NoError(t, err)
		}))
		defer server.Close()

		cms.url = server.URL
		response, err := cms.search(ctx, Parameters())
		assert.NoError(t, err)
		assert.Equal(t, 2, response.Total)
		assert.Equal(t, 0, response.Skip)
		assert.Equal(t, 100, response.Limit)
		assert.Equal(t, 2, len(response.Items))
	})

	t.Run("Should return ErrTooManyRequests if retry shouldn't be tried", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer server.Close()

		cms.url = server.URL

		t.Run("No deadline set on context", func(t *testing.T) {
			_, err := cms.search(context.Background(), Parameters())
			assert.Error(t, err)
			assert.Equal(t, ErrTooManyRequests, err)
		})

		t.Run("Deadline is before the retry seconds", func(t *testing.T) {
			c, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := cms.search(c, Parameters())
			assert.Error(t, err)
			assert.Equal(t, ErrTooManyRequests, err)
		})
	})

	t.Run("Should retry if hits Contentful rate limit", func(t *testing.T) {
		var (
			rateLimit       = true
			happyCaseCalled = false
		)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rateLimit {
				w.Header().Set("X-Contentful-RateLimit-Reset", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				rateLimit = false
				return
			}

			happyCaseCalled = true

			w.WriteHeader(http.StatusOK)
			bytes, err := ioutil.ReadFile("testdata/prod_all_pages.json")
			assert.NoError(t, err)
			_, err = w.Write(bytes)
			assert.NoError(t, err)
		}))
		defer server.Close()

		cms.url = server.URL

		c, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := cms.search(c, Parameters())
		assert.NoError(t, err)
		assert.True(t, happyCaseCalled)
	})
}

func TestContentful_retryAfter(t *testing.T) {
	t.Parallel()

	t.Run("If deadline is not set on context, will always return -1", func(t *testing.T) {
		ctx := context.TODO()
		assert.Equal(t, -1, retryAfter(ctx, nil))

		ctx = context.Background()
		assert.Equal(t, -1, retryAfter(ctx, nil))

		ctx = context.WithValue(context.Background(), struct{}{}, "bar")
		assert.Equal(t, -1, retryAfter(ctx, nil))
	})

	t.Run("If deadline is set on context, will not return -1", func(t *testing.T) {
		ctx1, cancel1 := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
		defer cancel1()
		assert.NotEqual(t, -1, retryAfter(ctx1, nil))

		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Hour*1)
		defer cancel2()
		assert.NotEqual(t, -1, retryAfter(ctx2, nil))
	})

	t.Run("If header can't be parsed, returns 2", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
		defer cancel()
		assert.Equal(t, 2, retryAfter(ctx, nil))

		resp := &http.Response{
			Header: make(http.Header),
		}
		assert.Equal(t, 2, retryAfter(ctx, resp))

		resp.Header.Set("X-Contentful-RateLimit-Reset", "foo")
		assert.Equal(t, 2, retryAfter(ctx, resp))

		resp.Header.Set("X-Contentful-RateLimit-Reset", "5")
		assert.NotEqual(t, 2, retryAfter(ctx, resp))
	})

	t.Run("If deadline of context would end before retry, will return -1", func(t *testing.T) {
		t.Run("Header not ok", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			defer cancel()

			assert.Equal(t, -1, retryAfter(ctx, nil))
		})

		t.Run("Header ok", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			resp := &http.Response{
				Header: make(http.Header),
			}
			resp.Header.Set("X-Contentful-RateLimit-Reset", "10")
			assert.Equal(t, -1, retryAfter(ctx, resp))
		})
	})

	t.Run("If deadline of context will end after retry, will return retry seconds", func(t *testing.T) {
		t.Run("Header not ok", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			assert.Equal(t, 2, retryAfter(ctx, nil))
		})

		t.Run("Header not ok", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp := &http.Response{
				Header: make(http.Header),
			}
			resp.Header.Set("X-Contentful-RateLimit-Reset", "5")
			assert.Equal(t, 5, retryAfter(ctx, resp))
		})
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
