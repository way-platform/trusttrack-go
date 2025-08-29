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

// GetObjectGroupRequest is the request for the [Client.GetObjectGroup] method.
type GetObjectGroupRequest struct {
	// The external ID of the object group.
	ExternalID string `json:"externalId"`
}

// Query returns the query parameters for the request.
func (r *GetObjectGroupRequest) Query() url.Values {
	q := url.Values{}
	q.Set("version", "1")
	return q
}

// GetObjectGroupResponse is the response for the [Client.GetObjectGroup] method.
type GetObjectGroupResponse struct {
	// The object group.
	ObjectGroup *trusttrackv1.ObjectGroup `json:"objectGroup"`
}

// GetObjectGroup gets a specific object group by external ID.
func (c *Client) GetObjectGroup(
	ctx context.Context,
	request *GetObjectGroupRequest,
) (*GetObjectGroupResponse, error) {
	httpResponse, err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/object-groups/%s", url.PathEscape(request.ExternalID)),
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
	var responseBody ttoapi.ExternalObjectGroup
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	response := GetObjectGroupResponse{
		ObjectGroup: objectGroupToProto(&responseBody),
	}
	return &response, nil
}
