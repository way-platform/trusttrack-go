package trusttrack

import (
	"fmt"
	"io"
	"net/http"
)

// Error represents a LogiApp BI API error.
type Error struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Message is the error message from the response.
	Message string
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("http %d: %s", e.StatusCode, e.Message)
}

func newResponseError(response *http.Response) error {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	message := response.Status
	if len(data) > 0 {
		message = string(data)
	}
	return &Error{
		StatusCode: response.StatusCode,
		Message:    message,
	}
}
