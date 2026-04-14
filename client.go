package trusttrack

import (
	"net/http"
	"runtime/debug"
	"time"

	trusttrackv1connect "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1/trusttrackv1connect"
)

var _ trusttrackv1connect.TrustTrackApiClient = (*Client)(nil)

// Client for the TrustTrack Fleet Management API.
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
	baseURL        string
	apiKey         string
	baseHTTPClient *http.Client
	timeout        time.Duration
	retryCount     int
	interceptors   []func(http.RoundTripper) http.RoundTripper
}

func newClientConfig() clientConfig {
	return clientConfig{
		timeout:    10 * time.Second,
		retryCount: 3,
		baseURL:    "https://api.fm-track.com",
	}
}

func (cc clientConfig) httpClient() *http.Client {
	var transport http.RoundTripper
	if cc.baseHTTPClient != nil {
		transport = cc.baseHTTPClient.Transport
	}
	if transport == nil {
		transport = http.DefaultTransport
	}
	// Add API key transport if API key is configured.
	if cc.apiKey != "" {
		transport = &apiKeyTransport{
			apiKey:    cc.apiKey,
			transport: transport,
		}
	}
	// Add interceptor transport if interceptors are configured.
	if len(cc.interceptors) > 0 {
		transport = &interceptorTransport{
			interceptors: cc.interceptors,
			next:         transport,
		}
	}
	// Add retry transport if retry count > 0 (outermost).
	if cc.retryCount > 0 {
		transport = &retryTransport{
			next: transport,
		}
	}
	return &http.Client{
		Timeout:   cc.timeout,
		Transport: transport,
	}
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

// WithHTTPClient sets a custom HTTP client whose transport is used as the
// base for the SDK's transport chain.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *clientConfig) {
		c.baseHTTPClient = httpClient
	}
}

// WithInterceptor adds a request interceptor for the [Client].
func WithInterceptor(interceptor func(http.RoundTripper) http.RoundTripper) ClientOption {
	return func(c *clientConfig) {
		c.interceptors = append(c.interceptors, interceptor)
	}
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
