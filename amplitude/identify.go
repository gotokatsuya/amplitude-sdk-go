package amplitude

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Identification struct {
	UserID             string                 `json:"user_id,omitempty"`
	DeviceID           string                 `json:"device_id,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	AppVersion         string                 `json:"app_version,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	OSName             string                 `json:"os_name,omitempty"`
	OSVersion          string                 `json:"os_version,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	Country            string                 `json:"country,omitempty"`
	Region             string                 `json:"region,omitempty"`
	City               string                 `json:"city,omitempty"`
	DMA                string                 `json:"dma,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Paying             string                 `json:"paying,omitempty"`
	StartVersion       string                 `json:"start_version,omitempty"`
}

type IdentifyRequest struct {
	Identifications []Identification
}

func (r IdentifyRequest) toFormValues() (url.Values, error) {
	v := make(url.Values)

	b, err := json.Marshal(r.Identifications)
	if err != nil {
		return nil, err
	}

	v.Set("identification", string(b))
	return v, nil
}

// Identify ...
func (c *Client) Identify(ctx context.Context, req *IdentifyRequest) (*http.Response, error) {
	path := "/identify"
	values, err := req.toFormValues()
	if err != nil {
		return nil, err
	}
	httpReq, err := c.NewV1Request(http.MethodPost, path, values)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.Do(ctx, httpReq, nil)
	if err != nil {
		return httpResp, err
	}
	return httpResp, nil
}
