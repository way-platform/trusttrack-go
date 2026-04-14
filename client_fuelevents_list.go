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

// ListFuelEvents lists fuel events for an object.
func (c *Client) ListFuelEvents(
	ctx context.Context,
	request *trusttrackv1.ListFuelEventsRequest,
) (_ *trusttrackv1.ListFuelEventsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list fuel events: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "1")
	if request.GetObjectId() != "" {
		q.Set("object_id", request.GetObjectId())
	}
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
	fullURL := c.config.baseURL + "/fuel-events"
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
	var responseBody ttoapi.ExternalFuelEventsCollection
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.ListFuelEventsResponse{}
	fuelEvents := make([]*trusttrackv1.FuelEvent, 0, len(responseBody.Items))
	for _, fuelEvent := range responseBody.Items {
		fuelEvents = append(fuelEvents, fuelEventToProto(&fuelEvent))
	}
	resp.SetFuelEvents(fuelEvents)
	if responseBody.ContinuationToken != nil {
		resp.SetContinuationToken(strconv.Itoa(int(*responseBody.ContinuationToken)))
	}
	return resp, nil
}
