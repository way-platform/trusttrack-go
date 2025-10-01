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

// ListObjectsRequest is the request for the [Client.ListObjects] method.
type ListObjectsRequest struct{}

// Query returns the query parameters for the request.
func (r *ListObjectsRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "1")
	return q
}

// ListObjectsResponse is the response for the [Client.ListObjects] method.
type ListObjectsResponse struct {
	Objects []*trusttrackv1.Object `json:"objects"`
}

// ListObjects lists all objects.
func (c *Client) ListObjects(
	ctx context.Context,
	request *ListObjectsRequest,
	opts ...ClientOption,
) (_ *ListObjectsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list objects: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	fullURL := cfg.baseURL + "/objects"
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
