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

// ListObjectGroups lists all object groups.
func (c *Client) ListObjectGroups(
	ctx context.Context,
	request *trusttrackv1.ListObjectGroupsRequest,
) (_ *trusttrackv1.ListObjectGroupsResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list object groups: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "1")
	if request.GetLimit() > 0 {
		q.Set("limit", strconv.Itoa(int(request.GetLimit())))
	}
	if request.GetContinuationToken() != "" {
		q.Set("continuation_token", request.GetContinuationToken())
	}
	fullURL := c.config.baseURL + "/object-groups"
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
		Items             []ttoapi.ExternalObjectGroup `json:"items"`
		ContinuationToken *int32                       `json:"continuation_token,omitempty"`
	}
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.ListObjectGroupsResponse{}
	objectGroups := make([]*trusttrackv1.ObjectGroup, 0, len(responseBody.Items))
	for _, objectGroup := range responseBody.Items {
		objectGroups = append(objectGroups, objectGroupToProto(&objectGroup))
	}
	resp.SetObjectGroups(objectGroups)
	if responseBody.ContinuationToken != nil {
		resp.SetContinuationToken(fmt.Sprintf("%d", *responseBody.ContinuationToken))
	}
	return resp, nil
}
