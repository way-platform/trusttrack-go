package trusttrack

import "net/http"

// apiKeyTransport is an HTTP transport that authenticates requests using HTTP basic authentication.
type apiKeyTransport struct {
	apiKey    string
	transport http.RoundTripper
}

var _ http.RoundTripper = &apiKeyTransport{}

// RoundTrip implements the [http.RoundTripper] interface.
func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	query := req.URL.Query()
	query.Set("api_key", t.apiKey)
	req.URL.RawQuery = query.Encode()
	return t.transport.RoundTrip(req)
}
