package trusttrack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListObjectsLastCoordinateRequest is the request for the [Client.ListObjectsLastCoordinate] method.
type ListObjectsLastCoordinateRequest struct {
	// The limit of the number of objects to return.
	// Default: 100.
	// Maximum: 1000.
	Limit int `json:"limit"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken"`
}

// Query returns the query parameters for the request.
func (r *ListObjectsLastCoordinateRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "2")
	if r.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", r.Limit))
	}
	if r.ContinuationToken != "" {
		q.Set("continuation_token", r.ContinuationToken)
	}
	return q
}

// ListObjectsLastCoordinateResponse is the response for the [Client.ListObjectsLastCoordinate] method.
type ListObjectsLastCoordinateResponse struct {
	// The objects with their last coordinate.
	Objects []*trusttrackv1.Object `json:"objects"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken"`
}

// ListObjectsLastCoordinate lists all objects with their last coordinate.
func (c *Client) ListObjectsLastCoordinate(
	ctx context.Context,
	request *ListObjectsLastCoordinateRequest,
) (*ListObjectsLastCoordinateResponse, error) {
	httpResponse, err := c.doRequest(
		ctx,
		http.MethodGet,
		"/objects-last-coordinate",
		request.Query(),
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
	var responseBody struct {
		Results           []*ttoapi.ExternalComposedObject `json:"results"`
		ContinuationToken *string                          `json:"continuation_token"`
	}
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListObjectsLastCoordinateResponse{
		Objects: make([]*trusttrackv1.Object, 0, len(responseBody.Results)),
	}
	for _, object := range responseBody.Results {
		response.Objects = append(response.Objects, objectToProto(object))
	}
	if responseBody.ContinuationToken != nil {
		response.ContinuationToken = *responseBody.ContinuationToken
	}
	return &response, nil
}
