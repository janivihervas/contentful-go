package contentful

import (
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

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
		_, err := cms.search(ctx, nil)
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
		_, _ = cms.search(ctx, nil)
	})

	t.Run("Should return an error if Contentful responds with non-200 status code", func(tt *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cms.url = server.URL
		_, err := cms.search(ctx, nil)
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
		_, err := cms.search(ctx, nil)
		assert.NotNil(tt, err)
	})

	t.Run("Should not return an error and should return correctly parsed search results", func(tt *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bytes, err := ioutil.ReadFile("test_data/prod_all_pages.json")
			assert.Nil(tt, err)
			_, err = w.Write(bytes)
			assert.Nil(tt, err)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cms.url = server.URL
		response, err := cms.search(ctx, nil)
		assert.Nil(tt, err)
		assert.Equal(tt, 2, response.Total)
		assert.Equal(tt, 0, response.Skip)
		assert.Equal(tt, 100, response.Limit)
		assert.Equal(tt, 2, len(response.Items))
	})
}
