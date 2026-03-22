package trusttrack

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
)

// DebugTransport dumps HTTP requests and responses to stderr when enabled.
type DebugTransport struct {
	Enabled *bool
	Next    http.RoundTripper
}

func (t *DebugTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.Enabled == nil || !*t.Enabled {
		return t.next().RoundTrip(request)
	}
	requestDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump request for debug: %w", err)
	}
	prettyPrintDump(os.Stderr, requestDump, "> ")
	response, err := t.next().RoundTrip(request)
	if err != nil {
		return nil, err
	}
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump response for debug: %w", err)
	}
	prettyPrintDump(os.Stderr, responseDump, "< ")
	return response, nil
}

func (t *DebugTransport) next() http.RoundTripper {
	if t.Next != nil {
		return t.Next
	}
	return http.DefaultTransport
}

func prettyPrintDump(w io.Writer, dump []byte, prefix string) {
	var output bytes.Buffer
	output.Grow(len(dump) * 2)
	for line := range bytes.Lines(dump) {
		output.WriteString(prefix)
		output.Write(line)
	}
	output.WriteByte('\n')
	if _, err := w.Write(output.Bytes()); err != nil {
		slog.Warn("failed to write debug dump", "error", err)
	}
}
