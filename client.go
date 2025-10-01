package trusttrack

import (
	"net/http"
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
	baseURL      string
	apiKey       string
	debug        bool
	timeout      time.Duration
	retryCount   int
	interceptors []func(http.RoundTripper) http.RoundTripper
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

func (cc clientConfig) httpClient() *http.Client {
	transport := http.DefaultTransport
	// Add debug transport if debug is enabled (innermost).
	if cc.debug {
		transport = &debugTransport{
			next: transport,
		}
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

// WithDebug toggles debug mode (request/response dumps to stderr)..
func WithDebug(debug bool) ClientOption {
	return func(c *clientConfig) {
		c.debug = debug
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
