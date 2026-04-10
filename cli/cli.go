package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Credentials holds API credentials for the TrustTrack API.
type Credentials struct {
	APIKey string `json:"api_key"`
}

// CredentialStore loads, saves, and clears credentials.
type CredentialStore interface {
	Load() (*Credentials, error)
	Save(*Credentials) error
	Clear() error
}

// Option configures the CLI command tree.
type Option func(*config)

type config struct {
	credentialStore CredentialStore
	httpClient      *http.Client
}

// WithCredentialStore sets the credential store.
func WithCredentialStore(s CredentialStore) Option {
	return func(c *config) { c.credentialStore = s }
}

// WithHTTPClient sets a custom HTTP client for the SDK.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) { c.httpClient = httpClient }
}

// CredentialFileStore is a JSON file-backed credential store.
type CredentialFileStore struct {
	path string
}

// NewCredentialFileStore creates a new file-backed credential store at the given path.
func NewCredentialFileStore(path string) *CredentialFileStore {
	return &CredentialFileStore{path: path}
}

// Load reads credentials from the file.
func (s *CredentialFileStore) Load() (*Credentials, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("read store: %w", err)
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("unmarshal store: %w", err)
	}
	return &creds, nil
}

// Save writes credentials to the file.
func (s *CredentialFileStore) Save(creds *Credentials) error {
	out, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}
	return os.WriteFile(s.path, out, 0o600)
}

// Clear removes the credential file.
func (s *CredentialFileStore) Clear() error {
	err := os.Remove(s.path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}
