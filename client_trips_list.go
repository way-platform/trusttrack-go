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

// ListTrips lists trips for an object.
func (c *Client) ListTrips(
	ctx context.Context,
	request *trusttrackv1.ListTripsRequest,
) (_ *trusttrackv1.ListTripsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list trips: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "1")
	if request.HasFromTime() {
		q.Set("from_datetime", request.GetFromTime().AsTime().UTC().Format(time.RFC3339))
	}
	if request.HasToTime() {
		q.Set("to_datetime", request.GetToTime().AsTime().UTC().Format(time.RFC3339))
	}
	if request.GetLimit() > 0 {
		q.Set("limit", strconv.Itoa(int(request.GetLimit())))
	}
	if request.GetContinuationToken() != "" {
		q.Set("continuation_token", request.GetContinuationToken())
	}
	fullURL := c.config.baseURL + fmt.Sprintf("/objects/%s/trips", request.GetObjectId())
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.URL.RawQuery = q.Encode()
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := c.config.httpClient().Do(httpRequest)
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
	resp := &trusttrackv1.ListTripsResponse{}
	trips := make([]*trusttrackv1.Trip, 0, len(responseBody.Trips))
	for _, trip := range responseBody.Trips {
		trips = append(trips, tripToProto(&trip))
	}
	resp.SetTrips(trips)
	if responseBody.ContinuationToken != nil {
		resp.SetContinuationToken((*responseBody.ContinuationToken).Format(time.RFC3339))
	}
	return resp, nil
}
