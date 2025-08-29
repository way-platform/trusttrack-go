package trusttrack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

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
	httpRequest, err := c.newRequest(ctx, "GET", "/objects-last-coordinate?version=2", nil)
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.do(ctx, httpRequest)
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
