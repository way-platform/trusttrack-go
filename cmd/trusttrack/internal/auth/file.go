package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/way-platform/trusttrack-go"
)

// File storing authentication credentials for the CLI.
type File struct {
	// APIKey is API key to use for authentication.
	APIKey string `json:"username"`
}

// NewClient creates a new TrustTrack API client using the current CLI credentials.
func NewClient(opts ...trusttrack.ClientOption) (*trusttrack.Client, error) {
	cf, err := readFile()
	if err != nil {
		return nil, err
	}
	return trusttrack.NewClient(
		append(
			opts,
			trusttrack.WithAPIKey(cf.APIKey),
		)...,
	)
}

func resolveFilepath() (string, error) {
	return xdg.ConfigFile("trusttrack-go/auth.json")
}

// readFile reads the currently stored [File].
func readFile() (*File, error) {
	fp, err := resolveFilepath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no credentials found, please login using `trusttrack auth login`")
		}
		return nil, err
	}
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// writeFile writes the stored [File].
func writeFile(f *File) error {
	fp, err := resolveFilepath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fp, data, 0o600)
}

// removeFile removes the stored [File].
func removeFile() error {
	fp, err := resolveFilepath()
	if err != nil {
		return err
	}
	return os.RemoveAll(fp)
}
