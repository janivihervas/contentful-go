package contentful

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// GetMany entries from Contentful. The flattened json output will be marshaled into data parameter,
// which will need to be a slice or an array. Will return an error if zero entries were returned
func (cms *Contentful) GetMany(ctx context.Context, parameters SearchParameters, data interface{}) error {
	response, err := cms.search(ctx, parameters)
	if err != nil {
		return err
	}

	if response.Total == 0 || len(response.Items) == 0 {
		return errors.New("no items returned")
	}

	appendIncludes(&response)

	flattenedItems, err := flattenItems(response.Includes, response.Items)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(flattenedItems)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}

// GetOne entry from Contentful. The flattened json output will be marshaled into data parameter.
// Will return an error if there is not exactly one entry returned
func (cms *Contentful) GetOne(ctx context.Context, parameters SearchParameters, data interface{}) error {
	response, err := cms.search(ctx, parameters)
	if err != nil {
		return err
	}

	if response.Total != 1 || len(response.Items) != 1 {
		return fmt.Errorf("too many or too few items returned: %d", response.Total)
	}

	appendIncludes(&response)

	flattenedItem, err := flattenItem(response.Includes, response.Items[0])
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(flattenedItem)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}

func (cms *Contentful) search(ctx context.Context, parameters SearchParameters) (searchResults, error) {
	response := searchResults{}
	if parameters.Values == nil {
		parameters.Values = url.Values{}
	}
	parameters.Set("include", "10")

	u := cms.url + "/spaces/" + cms.spaceID + "/entries?" + parameters.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return response, err
	}

	req.Header.Add("Authorization", "Bearer "+cms.token)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("non-ok status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	return response, err
}

// appendIncludes will append current search results to includes object,
// because Contentful doesn't duplicate items from search results to includes.
func appendIncludes(response *searchResults) {
	for _, item := range response.Items {
		if item.Sys.Type == linkTypeEntry {
			response.Includes.Entry = append(response.Includes.Entry, item)
		}
	}
}
