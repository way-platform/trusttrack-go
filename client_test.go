package trusttrack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"connectrpc.com/connect"
	trusttrackv1 "github.com/way-platform/trusttrack-go/proto/gen/go/wayplatform/connect/trusttrack/v1"
)

// newTestClient creates a Client pointing at the given test server with retries disabled.
func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	client, err := NewClient(
		WithBaseURL(srv.URL),
		WithAPIKey("test-key"),
		WithRetryCount(0),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestListObjects_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/objects" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("api_key"); got != "test-key" {
			t.Errorf("expected api_key=test-key, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": "1", "name": "Truck A"},
			{"id": "2", "name": "Truck B"},
		})
	}))
	defer srv.Close()
	client := newTestClient(t, srv)
	resp, err := client.ListObjects(context.Background(), &trusttrackv1.ListObjectsRequest{})
	if err != nil {
		t.Fatalf("ListObjects: %v", err)
	}
	if got := len(resp.GetObjects()); got != 2 {
		t.Errorf("expected 2 objects, got %d", got)
	}
}

func TestListDrivers_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/drivers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("version"); got != "2" {
			t.Errorf("expected version=2, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "10", "first_name": "Driver", "last_name": "One"},
			},
		})
	}))
	defer srv.Close()
	client := newTestClient(t, srv)
	resp, err := client.ListDrivers(context.Background(), &trusttrackv1.ListDriversRequest{})
	if err != nil {
		t.Fatalf("ListDrivers: %v", err)
	}
	if got := len(resp.GetDrivers()); got != 1 {
		t.Errorf("expected 1 driver, got %d", got)
	}
}

func TestListObjectGroups_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/object-groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "100", "name": "Fleet Group"},
			},
		})
	}))
	defer srv.Close()
	client := newTestClient(t, srv)
	resp, err := client.ListObjectGroups(context.Background(), &trusttrackv1.ListObjectGroupsRequest{})
	if err != nil {
		t.Fatalf("ListObjectGroups: %v", err)
	}
	if got := len(resp.GetObjectGroups()); got != 1 {
		t.Errorf("expected 1 object group, got %d", got)
	}
}

func TestErrorCodeMapping(t *testing.T) {
	tests := []struct {
		httpStatus int
		wantCode   connect.Code
	}{
		{http.StatusBadRequest, connect.CodeInvalidArgument},
		{http.StatusUnauthorized, connect.CodeUnauthenticated},
		{http.StatusForbidden, connect.CodePermissionDenied},
		{http.StatusNotFound, connect.CodeNotFound},
		{http.StatusConflict, connect.CodeAlreadyExists},
		{http.StatusTooManyRequests, connect.CodeResourceExhausted},
		{http.StatusNotImplemented, connect.CodeUnimplemented},
		{http.StatusServiceUnavailable, connect.CodeUnavailable},
		{http.StatusGatewayTimeout, connect.CodeDeadlineExceeded},
		{http.StatusInternalServerError, connect.CodeInternal},
		{http.StatusTeapot, connect.CodeUnknown},
	}
	for _, tt := range tests {
		t.Run(http.StatusText(tt.httpStatus), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "something went wrong", tt.httpStatus)
			}))
			defer srv.Close()
			client := newTestClient(t, srv)
			_, err := client.ListObjects(context.Background(), &trusttrackv1.ListObjectsRequest{})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var connectErr *connect.Error
			if !errors.As(err, &connectErr) {
				t.Fatalf("expected connect error, got %T: %v", err, err)
			}
			if connectErr.Code() != tt.wantCode {
				t.Errorf("expected code %v, got %v", tt.wantCode, connectErr.Code())
			}
		})
	}
}

func TestInterceptor(t *testing.T) {
	var intercepted atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Error("interceptor did not set expected header")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]any{})
	}))
	defer srv.Close()
	client, err := NewClient(
		WithBaseURL(srv.URL),
		WithAPIKey("test-key"),
		WithRetryCount(0),
		WithInterceptor(func(next http.RoundTripper) http.RoundTripper {
			intercepted.Store(true)
			return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				req.Header.Set("X-Custom-Header", "test-value")
				return next.RoundTrip(req)
			})
		}),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, err = client.ListObjects(context.Background(), &trusttrackv1.ListObjectsRequest{})
	if err != nil {
		t.Fatalf("ListObjects: %v", err)
	}
	if !intercepted.Load() {
		t.Error("interceptor was not called")
	}
}

// roundTripperFunc adapts a function to http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
