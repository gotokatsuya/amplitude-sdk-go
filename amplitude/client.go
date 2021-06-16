package amplitude

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// NewV1Request method
func (c *Client) NewV1Request(method, path string, body url.Values) (*http.Request, error) {

	u, err := c.endpointV1.Parse(path)
	if err != nil {
		return nil, err
	}

	body.Set("api_key", c.apiKey)

	req, err := http.NewRequest(method, u.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// NewV2Request method
func (c *Client) NewV2Request(method, path string, body interface{}) (*http.Request, error) {

	u, err := c.endpointV2.Parse(path)
	if err != nil {
		return nil, err
	}

	if req, ok := body.(RequestV2); ok {
		req.SetAPIKey(c.apiKey)
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
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

type requestV2 struct {
	APIKey string `json:"api_key"`
}

func (r *requestV2) SetAPIKey(v string) {
	r.APIKey = v
}

type RequestV2 interface {
	SetAPIKey(string)
}

type responseV2 struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}
