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

// ListObjectGroupsRequest is the request for the [Client.ListObjectGroups] method.
type ListObjectGroupsRequest struct {
	// The limit of the number of object groups to return.
	// Default: 100.
	// Maximum: 1000.
	Limit int `json:"limit"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken"`
}

// Query returns the query parameters for the request.
func (r *ListObjectGroupsRequest) Query() url.Values {
	q := url.Values{}
	if r.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", r.Limit))
	}
	if r.ContinuationToken != "" {
		q.Set("continuation_token", r.ContinuationToken)
	}
	return q
}

// ListObjectGroupsResponse is the response for the [Client.ListObjectGroups] method.
type ListObjectGroupsResponse struct {
	// The object groups.
	ObjectGroups []*trusttrackv1.ObjectGroup `json:"objectGroups"`
	// The continuation token to use to get the next page of results.
	ContinuationToken string `json:"continuationToken"`
}

// ListObjectGroups lists all object groups.
func (c *Client) ListObjectGroups(
	ctx context.Context,
	request *ListObjectGroupsRequest,
) (*ListObjectGroupsResponse, error) {
	httpResponse, err := c.doRequest(
		ctx,
		http.MethodGet,
		"/object-groups",
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
		Items             []ttoapi.ExternalObjectGroup `json:"items"`
		ContinuationToken *int32                       `json:"continuation_token,omitempty"`
	}
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := ListObjectGroupsResponse{
		ObjectGroups: make([]*trusttrackv1.ObjectGroup, 0, len(responseBody.Items)),
	}
	for _, objectGroup := range responseBody.Items {
		response.ObjectGroups = append(response.ObjectGroups, objectGroupToProto(&objectGroup))
	}
	if responseBody.ContinuationToken != nil {
		response.ContinuationToken = fmt.Sprintf("%d", *responseBody.ContinuationToken)
	}
	return &response, nil
}
