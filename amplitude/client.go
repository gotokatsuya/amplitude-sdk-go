package amplitude

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// API endpoint base constants
const (
	APIEndpointV1 = "https://api.amplitude.com"
	APIEndpointV2 = "https://api2.amplitude.com"
)

// Client type
type Client struct {
	httpClient *http.Client
	endpointV1 *url.URL
	endpointV2 *url.URL

	apiKey string
}

// NewClient returns a new client instance.
func NewClient(apiKey string, httpClient *http.Client) (*Client, error) {
	c := &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
	u1, err := url.Parse(APIEndpointV1)
	if err != nil {
		return nil, err
	}
	c.endpointV1 = u1
	u2, err := url.Parse(APIEndpointV2)
	if err != nil {
		return nil, err
	}
	c.endpointV2 = u2
	return c, nil
}

// mergeQuery method
func (c *Client) mergeQuery(path string, q interface{}) (string, error) {
	v := reflect.ValueOf(q)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return path, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return path, err
	}

	qs, err := query.Values(q)
	if err != nil {
		return path, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewV1Request method
func (c *Client) NewV1Request(method, path string, body interface{}) (*http.Request, error) {
	return c.newRequest(c.endpointV1, method, path, body)
}

// NewV2Request method
func (c *Client) NewV2Request(method, path string, body interface{}) (*http.Request, error) {
	return c.newRequest(c.endpointV2, method, path, body)
}

// newRequest method
func (c *Client) newRequest(endpoint *url.URL, method, path string, body interface{}) (*http.Request, error) {

	if req, ok := body.(Request); ok {
		req.SetAPIKey(c.apiKey)
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		if body != nil {
			merged, err := c.mergeQuery(path, body)
			if err != nil {
				return nil, err
			}
			path = merged
		}
	}
	u, err := endpoint.Parse(path)
	if err != nil {
		return nil, err
	}

	var reqBody io.ReadWriter
	switch method {
	case http.MethodPost, http.MethodPut:
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(b)
		}
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Do method
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

type request struct {
	APIKey string `json:"api_key"`
}

func (r *request) SetAPIKey(v string) {
	r.APIKey = v
}

type Request interface {
	SetAPIKey(string)
}

type response struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}
