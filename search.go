package contentful

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.opencensus.io/trace"
)

var (
	// ErrNoEntries is returned if no entries were returned
	ErrNoEntries = errors.New("contentful: no entries returned")
	// ErrMoreThanOneEntry is returned if there were more than one entry returned
	ErrMoreThanOneEntry = errors.New("contentful: more then one entry was returned")
	// ErrTooManyRequests is returned if hit Contentful rate limit and context doesn't have a deadline set
	ErrTooManyRequests = errors.New("contentful: too many requests")
)

// GetMany entries from Contentful. The flattened json output will be marshaled into data parameter,
// which will need to be a slice or an array. Will return an error if zero entries were returned
//
// Will retry if Contentful rate limits the request if
// - context has a deadline/timeout set and
// - seconds to wait is not after context's deadline/timeout, making this fail early
func (cms *Contentful) GetMany(ctx context.Context, parameters SearchParameters, data interface{}) error {
	ctx, span := trace.StartSpan(ctx, "github.com/janivihervas/contentful-go.GetMany")
	defer span.End()

	response, err := cms.search(ctx, parameters)
	if err != nil {
		return err
	}

	if response.Total == 0 || len(response.Items) == 0 {
		addSpanError(span, trace.StatusCodeNotFound, ErrNoEntries)
		return ErrNoEntries
	}

	_, spanFlatten := trace.StartSpan(ctx, "github.com/janivihervas/contentful-go.flattenItems")
	appendIncludes(&response)

	flattenedItems, err := flattenItems(response.Includes, response.Items)
	if err != nil {
		addSpanError(spanFlatten, trace.StatusCodeUnknown, err)
		spanFlatten.End()
		return err
	}
	spanFlatten.End()

	bytes, err := json.Marshal(flattenedItems)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return err
	}

	return nil
}

// GetOne entry from Contentful. The flattened json output will be marshaled into data parameter.
// Will return an error if there is not exactly one entry returned
//
// Will retry if Contentful rate limits the request if
// - context has a deadline/timeout set and
// - seconds to wait is not after context's deadline/timeout, making this fail early
func (cms *Contentful) GetOne(ctx context.Context, parameters SearchParameters, data interface{}) error {
	ctx, span := trace.StartSpan(ctx, "github.com/janivihervas/contentful-go.GetOne")
	defer span.End()

	response, err := cms.search(ctx, parameters)
	if err != nil {
		return err
	}

	if response.Total == 0 || len(response.Items) == 0 {
		addSpanError(span, trace.StatusCodeNotFound, ErrNoEntries)
		return ErrNoEntries
	}

	if response.Total != 1 || len(response.Items) != 1 {
		addSpanError(span, trace.StatusCodeOutOfRange, ErrMoreThanOneEntry)
		return ErrMoreThanOneEntry
	}

	_, spanFlatten := trace.StartSpan(ctx, "github.com/janivihervas/contentful-go.flattenItem")
	appendIncludes(&response)

	flattenedItem, err := flattenItem(response.Includes, response.Items[0])
	if err != nil {
		addSpanError(spanFlatten, trace.StatusCodeUnknown, err)
		spanFlatten.End()
		return err
	}
	spanFlatten.End()

	bytes, err := json.Marshal(flattenedItem)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return err
	}

	return nil
}

func (cms *Contentful) search(ctx context.Context, parameters SearchParameters) (searchResults, error) {
	ctx, span := trace.StartSpan(ctx, "github.com/janivihervas/contentful-go.search")
	defer span.End()

	response := searchResults{}
	if parameters.Values == nil {
		parameters.Values = url.Values{}
	}
	parameters.Set("include", "10")

	urlStr := cms.url + "/spaces/" + cms.spaceID + "/entries?" + parameters.Encode()
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return response, err
	}

	span.AddAttributes(trace.StringAttribute("http.host", urlParsed.Host))
	span.AddAttributes(trace.StringAttribute("http.method", http.MethodGet))
	span.AddAttributes(trace.StringAttribute("http.path", urlParsed.Path))
	span.AddAttributes(trace.StringAttribute("http.query", urlParsed.RawQuery))

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return response, err
	}

	req.Header.Add("Authorization", "Bearer "+cms.token)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err == context.Canceled {
		addSpanError(span, trace.StatusCodeCancelled, err)
		return response, err
	}
	if err == context.DeadlineExceeded {
		addSpanError(span, trace.StatusCodeDeadlineExceeded, err)
		return response, err
	}
	if err != nil {
		addSpanError(span, trace.StatusCodeUnknown, err)
		return response, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	span.AddAttributes(trace.Int64Attribute("http.status_code", int64(resp.StatusCode)))

	if resp.StatusCode == http.StatusTooManyRequests {
		seconds := retryAfter(ctx, resp)
		if seconds == -1 {
			addSpanError(span, trace.StatusCodeCancelled, ErrTooManyRequests)
			return response, ErrTooManyRequests
		}

		span.AddAttributes(trace.Int64Attribute("http.ratelimit_reset", int64(seconds)))

		select {
		case <-time.After(time.Second * time.Duration(seconds)):
			addSpanError(span, trace.StatusCodeResourceExhausted, err)
			return cms.search(ctx, parameters)
		case <-ctx.Done():
			addSpanError(span, trace.StatusCodeCancelled, err)
			return response, ctx.Err()
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-ok status code: %d", resp.StatusCode)
		addSpanError(span, trace.StatusCodeUnknown, err)
		return response, err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		addSpanError(span, trace.StatusCodeInternal, err)
		return response, err
	}

	return response, nil
}

func retryAfter(ctx context.Context, resp *http.Response) int {
	timeUntilCancel, deadlineSet := ctx.Deadline()
	if !deadlineSet {
		return -1
	}

	var (
		retrySeconds = 2
		header       string
	)

	if resp != nil {
		header = resp.Header.Get("X-Contentful-RateLimit-Reset")
	}

	s, err := strconv.Atoi(header)
	if err == nil {
		retrySeconds = s
	}

	timeToRetry := time.Now().Add(time.Second * time.Duration(retrySeconds))
	shouldRetry := timeToRetry.Before(timeUntilCancel)
	if shouldRetry {
		return retrySeconds
	}

	return -1
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
