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

// ListObjects lists all objects.
func (c *Client) ListObjects(
	ctx context.Context,
	_ *trusttrackv1.ListObjectsRequest,
) (_ *trusttrackv1.ListObjectsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list objects: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "1")
	fullURL := c.config.baseURL + "/objects"
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
	var responseBody []*ttoapi.ExternalComposedObject
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.ListObjectsResponse{}
	objects := make([]*trusttrackv1.Object, 0, len(responseBody))
	for _, object := range responseBody {
		objects = append(objects, objectToProto(object))
	}
	resp.SetObjects(objects)
	return resp, nil
}
