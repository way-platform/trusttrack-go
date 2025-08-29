package trusttrack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListObjectCoordinatesRequest is the request for the [Client.ListObjectCoordinates] method.
type ListObjectCoordinatesRequest struct {
	// ObjectID is the external object ID.
	ObjectID string

	// FromTime finds records starting from the specified date and time.
	FromTime time.Time

	// ToTime finds records ending at the specified date and time. Optional.
	ToTime time.Time

	// ContinuationToken displays from what date and time the data is shown if the record limit was reached.
	ContinuationToken string

	// Limit specifies how many records will be included in the response (default 100, max 1000).
	Limit int

	// IncludeGeozones indicates whether to include geozone information in the response.
	IncludeGeozones bool

	// IncludeTireParameters indicates whether to include tire pressure information in the response.
	IncludeTireParameters bool
}

// Query returns the query parameters for the request.
func (r *ListObjectCoordinatesRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "2")
	if r.ObjectID != "" {
		q.Set("objectId", r.ObjectID)
	}
	if !r.FromTime.IsZero() {
		q.Set("from_datetime", r.FromTime.Format("2006-01-02T15:04:05.000Z"))
	}

	if !r.ToTime.IsZero() {
		q.Set("to_datetime", r.ToTime.Format("2006-01-02T15:04:05.000Z"))
	}
	if r.ContinuationToken != "" {
		q.Set("continuation_token", r.ContinuationToken)
	}
	if r.Limit > 0 {
		q.Set("limit", strconv.Itoa(r.Limit))
	}
	if r.IncludeGeozones {
		q.Set("include_geozones", "true")
	}
	if r.IncludeTireParameters {
		q.Set("include_tire_parameters", "true")
	}
	return q
}

// ListObjectCoordinatesResponse is the response for the [Client.ListObjectCoordinates] method.
type ListObjectCoordinatesResponse struct {
	Coordinates       []*trusttrackv1.Coordinate `json:"coordinates"`
	ContinuationToken string                     `json:"continuation_token"`
}

// ListObjectCoordinates lists object coordinates for a specified time period.
func (c *Client) ListObjectCoordinates(ctx context.Context, request *ListObjectCoordinatesRequest) (*ListObjectCoordinatesResponse, error) {
	requestPath := "/objects/" + request.ObjectID + "/coordinates"
	httpResponse, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestPath,
		request.Query(),
	)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
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
	response := ListObjectCoordinatesResponse{
		Coordinates: make([]*trusttrackv1.Coordinate, 0, len(responseBody.Items)),
	}
	if responseBody.ContinuationToken != nil {
		response.ContinuationToken = responseBody.ContinuationToken.Format(time.RFC3339)
	}
	for _, coordinate := range responseBody.Items {
		response.Coordinates = append(response.Coordinates, coordinateToProto(&coordinate))
	}
	return &response, nil
}
