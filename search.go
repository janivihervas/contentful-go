package contentful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Search entries from Contentful.
//
// Example usage:
//	type Page struct {
//		Title string `json:"title"`
//	}
//
//	cms := contentful.New("token", "space", false)
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
//	defer cancel()
//
//	page := Page{}
//  params := url.Values{
//    "content_type": []string{"page"}
//  }
//
//	err := cms.Search(ctx, params, &page)
//	if err != nil {
//	  // Handle error
//	}
func (cms *Contentful) Search(ctx context.Context, parameters url.Values, data interface{}) error {
	if parameters == nil {
		parameters = url.Values{}
	}
	parameters.Set("include", "10")

	u := cms.url + "/spaces/" + cms.spaceID + "/entries?" + parameters.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+cms.token)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-ok status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(data)
}
