package trusttrack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
)

// Client for the LogiApp BI API.
type Client struct {
	config     clientConfig
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new [Client] with the given options.
func NewClient(opts ...ClientOption) (*Client, error) {
	client := Client{
		httpClient: http.DefaultClient,
		config:     newClientConfig(),
	}
	for _, opt := range opts {
		opt(&client.config)
	}
	if client.config.apiKey == "" {
		return nil, fmt.Errorf("apiKey is required, use WithAPIKey() option")
	}
	client.baseURL = "https://api.fm-track.com"
	return &client, nil
}

// clientConfig is the config for a [Client].
type clientConfig struct {
	apiKey string
	debug  bool
}

func newClientConfig() clientConfig {
	return clientConfig{}
}

// ClientOption is a function that configures the [clientConfig].
type ClientOption func(*clientConfig)

// WithAPIKey sets the API key for API requests.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *clientConfig) {
		c.apiKey = apiKey
	}
}

// WithDebug toggles debug mode (request/response dumps to stderr)..
func WithDebug(debug bool) ClientOption {
	return func(c *clientConfig) {
		c.debug = debug
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	requestPath string,
	body io.Reader,
) (_ *http.Request, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("new request: %w", err)
		}
	}()
	fullURL := c.baseURL + requestPath
	request, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", getUserAgent())
	request.Header.Set("Accept", "application/json")
	return request, nil
}

func (c *Client) do(ctx context.Context, request *http.Request) (*http.Response, error) {
	query := request.URL.Query()
	query.Set("api_key", c.config.apiKey)
	request.URL.RawQuery = query.Encode()
	if c.config.debug {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s", dump)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if c.config.debug {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s", dump)
	}
	return response, nil
}

func getUserAgent() string {
	userAgent := "WayPlatformTrustTrackGo"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		userAgent += "/" + info.Main.Version
	}
	return userAgent
}
