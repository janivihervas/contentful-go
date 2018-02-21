package contentful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *client) Search(ctx context.Context, parameters url.Values, data interface{}) error {
	if parameters == nil {
		parameters = url.Values{}
	}
	parameters.Set("include", "10")

	u := c.url + "/spaces/" + c.spaceID + "/entries?" + parameters.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
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
