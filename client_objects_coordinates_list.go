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

// ListObjectCoordinates lists object coordinates for a specified time period.
func (c *Client) ListObjectCoordinates(
	ctx context.Context,
	request *trusttrackv1.ListObjectCoordinatesRequest,
) (_ *trusttrackv1.ListObjectCoordinatesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list object coordinates: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "2")
	if request.GetObjectId() != "" {
		q.Set("objectId", request.GetObjectId())
	}
	if request.HasFromTime() {
		q.Set("from_datetime", request.GetFromTime().AsTime().UTC().Format(time.RFC3339))
	}
	if request.HasToTime() {
		q.Set("to_datetime", request.GetToTime().AsTime().UTC().Format(time.RFC3339))
	}
	if request.GetContinuationToken() != "" {
		q.Set("continuation_token", request.GetContinuationToken())
	}
	if request.GetLimit() > 0 {
		q.Set("limit", strconv.Itoa(int(request.GetLimit())))
	} else {
		q.Set("limit", "1000")
	}
	if request.GetIncludeGeozones() {
		q.Set("include_geozones", "true")
	}
	if request.GetIncludeTireParameters() {
		q.Set("include_tire_parameters", "true")
	}
	requestPath := "/objects/" + request.GetObjectId() + "/coordinates"
	fullURL := c.config.baseURL + requestPath
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
	defer func() { _ = httpResponse.Body.Close() }()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var responseBody ttoapi.CoordinateCollection
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.ListObjectCoordinatesResponse{}
	coordinates := make([]*trusttrackv1.Coordinate, 0, len(responseBody.Items))
	for _, coordinate := range responseBody.Items {
		coordinates = append(coordinates, coordinateToProto(&coordinate))
	}
	resp.SetCoordinates(coordinates)
	if responseBody.ContinuationToken != nil {
		resp.SetContinuationToken(responseBody.ContinuationToken.Format(time.RFC3339))
	}
	return resp, nil
}
