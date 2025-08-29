package trusttrack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListObjectsRequest is the request for the [Client.ListObjects] method.
type ListObjectsRequest struct{}

// ListObjectsResponse is the response for the [Client.ListObjects] method.
type ListObjectsResponse struct {
	Objects []*trusttrackv1.Object `json:"objects"`
}

// ListObjects lists all objects.
func (c *Client) ListObjects(ctx context.Context, request *ListObjectsRequest) (*ListObjectsResponse, error) {
	httpRequest, err := c.newRequest(ctx, "GET", "/objects?version=1", nil)
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
	var responseBody []*ttoapi.ExternalComposedObject
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListObjectsResponse{
		Objects: make([]*trusttrackv1.Object, 0, len(responseBody)),
	}
	for _, object := range responseBody {
		response.Objects = append(response.Objects, objectToProto(object))
	}
	return &response, nil
}
