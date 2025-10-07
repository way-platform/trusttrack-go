package trusttrack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListDriversRequest is the request for the [Client.ListDrivers] method.
type ListDriversRequest struct {
	// The limit of the number of drivers to return.
	// Default: 100.
	// Maximum: 1000.
	Limit int `json:"limit"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken"`
	// The identifier type to filter by.
	IdentifierType string `json:"identifierType"`
	// The identifier to filter by.
	Identifier string `json:"identifier"`
}

// Query returns the query parameters for the request.
func (r *ListDriversRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "2")
	if r.Limit > 0 {
		q.Set("limit", strconv.Itoa(r.Limit))
	}
	if r.ContinuationToken != "" {
		q.Set("continuation_token", r.ContinuationToken)
	}
	if r.IdentifierType != "" {
		q.Set("identifier_type", r.IdentifierType)
	}
	if r.Identifier != "" {
		q.Set("identifier", r.Identifier)
	}
	return q
}

// ListDriversResponse is the response for the [Client.ListDrivers] method.
type ListDriversResponse struct {
	// The drivers.
	Drivers []*trusttrackv1.Driver `json:"drivers"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken"`
}

// ListDrivers lists all drivers.
func (c *Client) ListDrivers(
	ctx context.Context,
	request *ListDriversRequest,
	opts ...ClientOption,
) (_ *ListDriversResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list drivers: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	fullURL := cfg.baseURL + "/drivers"
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.URL.RawQuery = request.Query().Encode()
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := cfg.httpClient().Do(httpRequest)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var responseBody ttoapi.ExternalDriverContinuableList
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListDriversResponse{
		Drivers: make([]*trusttrackv1.Driver, 0, len(responseBody.Items)),
	}
	for _, driver := range responseBody.Items {
		response.Drivers = append(response.Drivers, driverToProto(&driver))
	}
	if responseBody.ContinuationToken != nil {
		response.ContinuationToken = strconv.Itoa(*responseBody.ContinuationToken)
	}
	return &response, nil
}
