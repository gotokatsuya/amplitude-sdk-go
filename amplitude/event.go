package amplitude

import (
	"context"
	"net/http"
)

type Event struct {
	UserID              string                 `json:"user_id,omitempty"`
	DeviceID            string                 `json:"device_id,omitempty"`
	EventType           string                 `json:"event_type"`
	Time                int64                  `json:"time,omitempty"`
	EventProperties     map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties      map[string]interface{} `json:"user_properties,omitempty"`
	Groups              map[string]interface{} `json:"groups,omitempty"`
	GroupProperties     map[string]interface{} `json:"group_properties,omitempty"`
	AppVersion          string                 `json:"app_version,omitempty"`
	Platform            string                 `json:"platform,omitempty"`
	OSName              string                 `json:"os_name,omitempty"`
	OSVersion           string                 `json:"os_version,omitempty"`
	DeviceBrand         string                 `json:"device_brand,omitempty"`
	DeviceManufacturer  string                 `json:"device_manufacturer,omitempty"`
	DeviceModel         string                 `json:"device_model,omitempty"`
	Carrier             string                 `json:"carrier,omitempty"`
	Country             string                 `json:"country,omitempty"`
	Region              string                 `json:"region,omitempty"`
	City                string                 `json:"city,omitempty"`
	DMA                 string                 `json:"dma,omitempty"`
	Language            string                 `json:"language,omitempty"`
	Price               float64                `json:"price,omitempty"`
	Quantity            int                    `json:"quantity,omitempty"`
	Revenue             float64                `json:"revenue,omitempty"`
	ProductID           string                 `json:"productId,omitempty"`
	RevenueType         string                 `json:"revenueType,omitempty"`
	LocationLat         float64                `json:"location_lat,omitempty"`
	LocationLng         float64                `json:"location_lng,omitempty"`
	IP                  string                 `json:"ip,omitempty"`
	IOSAdvertiserID     string                 `json:"idfa,omitempty"`
	IOSVendorID         string                 `json:"idfv,omitempty"`
	AndroidAdvertiserID string                 `json:"adid,omitempty"`
	AndroidID           string                 `json:"android_id,omitempty"`
	EventID             int                    `json:"event_id,omitempty"`
	SessionID           int                    `json:"session_id,omitempty"`
	InsertID            string                 `json:"insert_id,omitempty"`
}

type LogEventRequest struct {
	requestV2
	Events []Event `json:"events"`
}

type LogEventResponse struct {
	responseV2
	EventsIngested   int   `json:"events_ingested,omitempty"`
	PayloadSizeBytes int   `json:"payload_size_bytes,omitempty"`
	ServerUploadTime int64 `json:"server_upload_time,omitempty"`
}

// LogEvent ...
func (c *Client) LogEvent(ctx context.Context, req *LogEventRequest) (*LogEventResponse, *http.Response, error) {
	path := "/2/httpapi"
	httpReq, err := c.NewV2Request(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(LogEventResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}
