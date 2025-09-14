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

// ListTripsRequest is the request for the [Client.ListTrips] method.
type ListTripsRequest struct {
	// The ID of the object to get trips for.
	ObjectID string `json:"objectId"`
	// The start time for the trip search.
	FromTime time.Time `json:"fromDatetime"`
	// The end time for the trip search (optional).
	ToTime time.Time `json:"toDatetime,omitzero"`
	// The limit of the number of trips to return.
	// Default: 100.
	Limit int `json:"limit"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken,omitempty"`
}

// Query returns the query parameters for the request.
func (r *ListTripsRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "1")
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

// ListTripsResponse is the response for the [Client.ListTrips] method.
type ListTripsResponse struct {
	// The trips for the object.
	Trips []*trusttrackv1.Trip `json:"trips"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken,omitempty"`
}

// ListTrips lists trips for an object.
func (c *Client) ListTrips(
	ctx context.Context,
	request *ListTripsRequest,
	opts ...ClientOption,
) (_ *ListTripsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list trips: %w", err)
		}
	}()
	httpResponse, err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/objects/%s/trips", request.ObjectID),
		request.Query(),
		opts...,
	)
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
	var responseBody ttoapi.TripCollection
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListTripsResponse{
		Trips: make([]*trusttrackv1.Trip, 0, len(responseBody.Trips)),
	}
	for _, trip := range responseBody.Trips {
		response.Trips = append(response.Trips, tripToProto(&trip))
	}
	if responseBody.ContinuationToken != nil {
		response.ContinuationToken = (*responseBody.ContinuationToken).Format(time.RFC3339)
	}
	return &response, nil
}
