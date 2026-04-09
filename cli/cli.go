package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Store reads and writes JSON-serializable data.
type Store interface {
	Read(target any) error
	Write(data any) error
	Clear() error
}

// Option configures the CLI command tree.
type Option func(*config)

type config struct {
	credentialStore Store
	httpClient      *http.Client
}

// WithCredentialStore sets the credential store.
func WithCredentialStore(s Store) Option {
	return func(c *config) { c.credentialStore = s }
}

// WithHTTPClient sets a custom HTTP client for the SDK.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) { c.httpClient = httpClient }
}

// FileStore is a JSON file-backed store.
type FileStore struct {
	path string
}

// NewFileStore creates a new file-backed store at the given path.
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

// Read unmarshals the file contents into target. If target implements
// proto.Message, protojson is used; otherwise standard JSON.
func (s *FileStore) Read(target any) error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read store: %w", err)
	}
	if msg, ok := target.(proto.Message); ok {
		if err := protojson.Unmarshal(data, msg); err != nil {
			return fmt.Errorf("unmarshal store: %w", err)
		}
		return nil
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("unmarshal store: %w", err)
	}
	return nil
}

// Write marshals data and writes it to the file. If data implements
// proto.Message, protojson is used; otherwise standard JSON.
func (s *FileStore) Write(data any) error {
	var out []byte
	var err error
	if msg, ok := data.(proto.Message); ok {
		out, err = protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(msg)
	} else {
		out, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}
	return os.WriteFile(s.path, out, 0o600)
}

// Clear removes the file.
func (s *FileStore) Clear() error {
	err := os.Remove(s.path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}
