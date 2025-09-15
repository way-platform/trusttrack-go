package trusttrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

// Client for the LogiApp BI API.
type Client struct {
	config clientConfig
}

// NewClient creates a new [Client] with the given options.
func NewClient(opts ...ClientOption) (*Client, error) {
	client := Client{
		config: newClientConfig(),
	}
	for _, opt := range opts {
		opt(&client.config)
	}
	return &client, nil
}

// clientConfig is the config for a [Client].
type clientConfig struct {
	baseURL    string
	apiKey     string
	debug      bool
	timeout    time.Duration
	retryCount int
}

func newClientConfig() clientConfig {
	return clientConfig{
		timeout:    10 * time.Second,
		retryCount: 3,
		baseURL:    "https://api.fm-track.com",
	}
}

func (cc clientConfig) with(opts ...ClientOption) clientConfig {
	for _, opt := range opts {
		opt(&cc)
	}
	return cc
}

func (cc clientConfig) client() *http.Client {
	transport := http.DefaultTransport
	if cc.retryCount > 0 {
		transport = &retryTransport{next: transport}
	}
	if cc.debug {
		transport = &debugTransport{next: transport}
	}
	rt := &http.Client{
		Timeout:   cc.timeout,
		Transport: transport,
	}
	return rt
}

// ClientOption is a function that configures the [clientConfig].
type ClientOption func(*clientConfig)

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

// WithRetryCount sets the number of retries for API requests.
func WithRetryCount(retryCount int) ClientOption {
	return func(c *clientConfig) {
		c.retryCount = retryCount
	}
}

// WithAPIKey sets the API key for API requests.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *clientConfig) {
		c.apiKey = apiKey
	}
}

// WithTimeout sets the timeout for API requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) {
		c.timeout = timeout
	}
}

// WithDebug toggles debug mode (request/response dumps to stderr)..
func WithDebug(debug bool) ClientOption {
	return func(c *clientConfig) {
		c.debug = debug
	}
}

func (c *Client) doRequest(
	ctx context.Context,
	method string,
	requestPath string,
	query url.Values,
	opts ...ClientOption,
) (_ *http.Response, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("http request: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	fullURL := cfg.baseURL + requestPath
	request, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, err
	}
	if cfg.apiKey != "" {
		query.Set("api_key", cfg.apiKey)
	}
	request.URL.RawQuery = query.Encode()
	request.Header.Set("User-Agent", getUserAgent())
	request.Header.Set("Accept", "application/json")
	return cfg.client().Do(request)
}

func getUserAgent() string {
	userAgent := "WayPlatformTrustTrackGo"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/way-platform/trusttrack-go" {
				if dep.Version != "" && dep.Version != "v0.0.0-00010101000000-000000000000" {
					userAgent += "/" + dep.Version
				}
				break
			}
		}
	}
	return userAgent
}
