package contentful

import (
	"context"
	"fmt"
	"net/url"
)

func (c *client) Search(ctx context.Context, parameters url.Values, data interface{}) error {
	if parameters == nil {
		parameters = url.Values{}
	}
	parameters.Add("include", "10")

	fmt.Println("Params:", parameters.Encode())
	return nil
}
