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

// ListDrivers lists all drivers.
func (c *Client) ListDrivers(
	ctx context.Context,
	request *trusttrackv1.ListDriversRequest,
) (_ *trusttrackv1.ListDriversResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("trusttrack: list drivers: %w", err)
		}
	}()
	q := url.Values{}
	q.Set("version", "2")
	if request.GetLimit() > 0 {
		q.Set("limit", strconv.Itoa(int(request.GetLimit())))
	}
	if request.GetContinuationToken() != "" {
		q.Set("continuation_token", request.GetContinuationToken())
	}
	if request.GetIdentifierType() != "" {
		q.Set("identifier_type", request.GetIdentifierType())
	}
	if request.GetIdentifier() != "" {
		q.Set("identifier", request.GetIdentifier())
	}
	fullURL := c.config.baseURL + "/drivers"
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
	var responseBody ttoapi.ExternalDriverContinuableList
	if err := json.Unmarshal(responseData, &responseBody); err != nil {
		return nil, err
	}
	resp := &trusttrackv1.ListDriversResponse{}
	drivers := make([]*trusttrackv1.Driver, 0, len(responseBody.Items))
	for _, driver := range responseBody.Items {
		drivers = append(drivers, driverToProto(&driver))
	}
	resp.SetDrivers(drivers)
	if responseBody.ContinuationToken != nil {
		resp.SetContinuationToken(strconv.Itoa(*responseBody.ContinuationToken))
	}
	return resp, nil
}
