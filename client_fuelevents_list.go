package trusttrack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListFuelEventsRequest is the request for the [Client.ListFuelEvents] method.
type ListFuelEventsRequest struct {
	// The ID of the object to get fuel events for.
	ObjectID string `json:"objectId"`
	// The start time for the fuel events search.
	FromTime time.Time `json:"fromDatetime"`
	// The end time for the fuel events search (optional).
	ToTime time.Time `json:"toDatetime,omitzero"`
	// The limit of the number of fuel events to return.
	// Default: 100.
	Limit int `json:"limit"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken,omitempty"`
}

// Query returns the query parameters for the request.
func (r *ListFuelEventsRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "1")
	if r.ObjectID != "" {
		q.Set("object_id", r.ObjectID)
	}
	if !r.FromTime.IsZero() {
		q.Set("from_datetime", r.FromTime.UTC().Format(time.RFC3339))
	}
	if !r.ToTime.IsZero() {
		q.Set("to_datetime", r.ToTime.UTC().Format(time.RFC3339))
	}
	if r.Limit > 0 {
		q.Set("limit", strconv.Itoa(r.Limit))
	}
	if r.ContinuationToken != "" {
		q.Set("continuation_token", r.ContinuationToken)
	}
	return q
}

// ListFuelEventsResponse is the response for the [Client.ListFuelEvents] method.
type ListFuelEventsResponse struct {
	// The fuel events for the object.
	FuelEvents []*trusttrackv1.FuelEvent `json:"fuelEvents"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken,omitempty"`
}

// ListFuelEvents lists fuel events for an object.
func (c *Client) ListFuelEvents(
	ctx context.Context,
	request *ListFuelEventsRequest,
	opts ...ClientOption,
) (_ *ListFuelEventsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list fuel events: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	fullURL := cfg.baseURL + "/fuel-events"
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
	var responseBody ttoapi.ExternalFuelEventsCollection
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListFuelEventsResponse{
		FuelEvents: make([]*trusttrackv1.FuelEvent, 0, len(responseBody.Items)),
	}
	for _, fuelEvent := range responseBody.Items {
		response.FuelEvents = append(response.FuelEvents, fuelEventToProto(&fuelEvent))
	}
	if responseBody.ContinuationToken != nil {
		response.ContinuationToken = strconv.Itoa(int(*responseBody.ContinuationToken))
	}
	return &response, nil
}
