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

// GetObjectGroup gets a specific object group by external ID.
func (c *Client) GetObjectGroup(
	ctx context.Context,
	request *trusttrackv1.GetObjectGroupRequest,
) (_ *trusttrackv1.GetObjectGroupResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: get object group: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "1")
	fullURL := c.config.baseURL + fmt.Sprintf("/object-groups/%s", url.PathEscape(request.GetExternalId()))
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
	var responseBody ttoapi.ExternalObjectGroup
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.GetObjectGroupResponse{}
	resp.SetObjectGroup(objectGroupToProto(&responseBody))
	return resp, nil
}
