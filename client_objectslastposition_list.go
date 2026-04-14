package trusttrack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/way-platform/trusttrack-go/internal/oapi/ttoapi"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// ListObjectsLastPosition lists all objects with their last position.
func (c *Client) ListObjectsLastPosition(
	ctx context.Context,
	request *trusttrackv1.ListObjectsLastPositionRequest,
) (_ *trusttrackv1.ListObjectsLastPositionResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list objects last position: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "2")
	if request.GetLimit() > 0 {
		q.Set("limit", strconv.Itoa(int(request.GetLimit())))
	} else {
		q.Set("limit", "1000")
	}
	if request.GetContinuationToken() != "" {
		q.Set("continuation_token", request.GetContinuationToken())
	}
	fullURL := c.config.baseURL + "/objects-last-coordinate"
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
	var responseBody struct {
		Results           []*ttoapi.ExternalComposedObject `json:"results"`
		ContinuationToken *string                          `json:"continuation_token"`
	}
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.ListObjectsLastPositionResponse{}
	objects := make([]*trusttrackv1.Object, 0, len(responseBody.Results))
	for _, object := range responseBody.Results {
		objects = append(objects, objectToProto(object))
	}
	resp.SetObjects(objects)
	if responseBody.ContinuationToken != nil {
		resp.SetContinuationToken(*responseBody.ContinuationToken)
	}
	return resp, nil
}
