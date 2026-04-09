package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"buf.build/go/protovalidate"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/trusttrack-go"
	"golang.org/x/term"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// resolveCredentials returns credentials from the first available source:
// credential provider, then credential store.
func resolveCredentials(cfg *config) (*Credentials, error) {
	if cfg.credentialProvider != nil {
		creds, err := cfg.credentialProvider()
		if err != nil {
			return nil, fmt.Errorf("credential provider: %w", err)
		}
		if creds != nil {
			return creds, nil
		}
	}
	if cfg.credentialStore != nil {
		var creds Credentials
		if err := cfg.credentialStore.Read(&creds); err != nil {
			return nil, err
		}
		return &creds, nil
	}
	return nil, fmt.Errorf("no credential source configured")
}

// NewCommand builds the full CLI command tree for the TrustTrack SDK.
func NewCommand(opts ...Option) *cobra.Command {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	cmd := &cobra.Command{
		Use:   "trusttrack",
		Short: "TrustTrack API CLI",
	}
	cmd.PersistentFlags().Bool("validate", false, "Enable message validation")
	cmd.AddGroup(&cobra.Group{ID: "objects", Title: "Objects"})
	cmd.AddCommand(newListObjectsCommand(&cfg))
	cmd.AddCommand(newListObjectsLastPositionCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "object-groups", Title: "Object Groups"})
	cmd.AddCommand(newListObjectGroupsCommand(&cfg))
	cmd.AddCommand(newGetObjectGroupCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "drivers", Title: "Drivers"})
	cmd.AddCommand(newListDriversCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "coordinates", Title: "Coordinates"})
	cmd.AddCommand(newListObjectCoordinatesCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "trips", Title: "Trips"})
	cmd.AddCommand(newListTripsCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "fuel-events", Title: "Fuel Events"})
	cmd.AddCommand(newListFuelEventsCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "auth", Title: "Authentication"})
	cmd.AddCommand(newAuthCommand(&cfg))
	cmd.AddGroup(&cobra.Group{ID: "utils", Title: "Utils"})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

func newAuthCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Authenticate to the TrustTrack API",
		GroupID: "auth",
	}
	cmd.AddCommand(newLoginCommand(cfg))
	cmd.AddCommand(newLogoutCommand(cfg))
	return cmd
}

func newLoginCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the TrustTrack API",
	}
	apiKey := cmd.Flags().String("api-key", "", "API key for authentication")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		// Try credential provider first, then stored credentials.
		creds := &Credentials{}
		if cfg.credentialProvider != nil {
			provided, err := cfg.credentialProvider()
			if err != nil {
				return fmt.Errorf("credential provider: %w", err)
			}
			if provided != nil {
				creds = provided
			}
		}
		if creds.APIKey == "" && cfg.credentialStore != nil {
			if err := cfg.credentialStore.Read(creds); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("read credentials: %w", err)
			}
		}
		// Override with flag.
		if *apiKey != "" {
			creds.APIKey = *apiKey
		}
		// Prompt for missing API key.
		if creds.APIKey == "" {
			val, err := promptSecret(cmd, "Enter API key: ")
			if err != nil {
				return err
			}
			creds.APIKey = val
		}
		// Persist credentials.
		if cfg.credentialStore != nil {
			if err := cfg.credentialStore.Write(creds); err != nil {
				return fmt.Errorf("write credentials: %w", err)
			}
		}
		cmd.Println("Logged in to the TrustTrack API.")
		return nil
	}
	return cmd
}

func newLogoutCommand(cfg *config) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from the TrustTrack API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cfg.credentialStore != nil {
				if err := cfg.credentialStore.Clear(); err != nil {
					return fmt.Errorf("clear credentials: %w", err)
				}
			}
			cmd.Println("Logged out.")
			return nil
		},
	}
}

func newClient(_ *cobra.Command, cfg *config) (*trusttrack.Client, error) {
	creds, err := resolveCredentials(cfg)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("no credentials found, please login using `trusttrack auth login`")
		}
		return nil, fmt.Errorf("resolve credentials: %w", err)
	}
	opts := []trusttrack.ClientOption{
		trusttrack.WithAPIKey(creds.APIKey),
	}
	if cfg.httpClient != nil {
		opts = append(opts, trusttrack.WithHTTPClient(cfg.httpClient))
	}
	return trusttrack.NewClient(opts...)
}

func newListObjectsCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects",
		Short:   "List objects",
		GroupID: "objects",
	}
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		response, err := client.ListObjects(cmd.Context(), &trusttrack.ListObjectsRequest{})
		if err != nil {
			return err
		}
		for _, object := range response.Objects {
			printJSON(cmd, object)
			validate(cmd, object)
		}
		return nil
	}
	return cmd
}

func newListObjectsLastPositionCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects-last-position",
		Short:   "List objects with their last position",
		GroupID: "objects",
	}
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.ListObjectsLastPositionRequest{
			Limit: 1000,
		}
		for {
			response, err := client.ListObjectsLastPosition(cmd.Context(), &request)
			if err != nil {
				return err
			}
			for _, object := range response.Objects {
				printJSON(cmd, object)
				validate(cmd, object)
			}
			request.ContinuationToken = response.ContinuationToken
			if request.ContinuationToken == "" {
				break
			}
		}
		return nil
	}
	return cmd
}

func newListObjectCoordinatesCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "coordinates [object-id]",
		Short:   "List object coordinates for a time period",
		GroupID: "coordinates",
		Args:    cobra.ExactArgs(1),
	}
	fromTime := cmd.Flags().Time(
		"from",
		time.Now().Add(-24*time.Hour),
		[]string{time.DateOnly, time.RFC3339},
		"From time",
	)
	toTime := cmd.Flags().Time(
		"to",
		time.Now(),
		[]string{time.DateOnly, time.RFC3339},
		"To time",
	)
	includeGeozones := cmd.Flags().Bool("include-geozones", false, "Include geozone information")
	includeTireParameters := cmd.Flags().Bool("include-tire-parameters", false, "Include tire pressure information")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.ListObjectCoordinatesRequest{
			ObjectID:              args[0],
			FromTime:              *fromTime,
			ToTime:                *toTime,
			Limit:                 1000,
			IncludeGeozones:       *includeGeozones,
			IncludeTireParameters: *includeTireParameters,
		}
		for {
			response, err := client.ListObjectCoordinates(cmd.Context(), &request)
			if err != nil {
				return err
			}
			for _, coordinate := range response.Coordinates {
				printJSON(cmd, coordinate)
				validate(cmd, coordinate)
			}
			if response.ContinuationToken == "" {
				break
			}
			request.ContinuationToken = response.ContinuationToken
		}
		return nil
	}
	return cmd
}

func newListObjectGroupsCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object-groups",
		Short:   "List object groups",
		GroupID: "object-groups",
	}
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.ListObjectGroupsRequest{
			Limit: 1000,
		}
		for {
			response, err := client.ListObjectGroups(cmd.Context(), &request)
			if err != nil {
				return err
			}
			for _, objectGroup := range response.ObjectGroups {
				printJSON(cmd, objectGroup)
				validate(cmd, objectGroup)
			}
			request.ContinuationToken = response.ContinuationToken
			if request.ContinuationToken == "" {
				break
			}
		}
		return nil
	}
	return cmd
}

func newGetObjectGroupCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object-group [external-id]",
		Short:   "Get a specific object group by external ID",
		GroupID: "object-groups",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.GetObjectGroupRequest{
			ExternalID: args[0],
		}
		response, err := client.GetObjectGroup(cmd.Context(), &request)
		if err != nil {
			return err
		}
		printJSON(cmd, response.ObjectGroup)
		validate(cmd, response.ObjectGroup)
		return nil
	}
	return cmd
}

func newListDriversCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "drivers",
		Short:   "List drivers",
		GroupID: "drivers",
	}
	identifierType := cmd.Flags().String("identifier-type", "", "Filter by identifier type")
	identifier := cmd.Flags().String("identifier", "", "Filter by identifier value")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if *identifier != "" && *identifierType == "" {
			return fmt.Errorf("identifier-type is required when identifier is provided")
		}
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.ListDriversRequest{
			Limit:          1000,
			IdentifierType: *identifierType,
			Identifier:     *identifier,
		}
		for {
			response, err := client.ListDrivers(cmd.Context(), &request)
			if err != nil {
				return err
			}
			for _, driver := range response.Drivers {
				printJSON(cmd, driver)
				validate(cmd, driver)
			}
			if response.ContinuationToken == "" {
				break
			}
			request.ContinuationToken = response.ContinuationToken
		}
		return nil
	}
	return cmd
}

func newListTripsCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trips [object-id]",
		Short:   "List trips for an object",
		GroupID: "trips",
		Args:    cobra.ExactArgs(1),
	}
	fromTime := cmd.Flags().Time("from", time.Now().Add(-24*time.Hour), []string{time.RFC3339}, "From time (RFC3339 format)")
	toTime := cmd.Flags().Time("to", time.Time{}, []string{time.RFC3339}, "To time (RFC3339 format, optional)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.ListTripsRequest{
			ObjectID: args[0],
			FromTime: *fromTime,
			ToTime:   *toTime,
			Limit:    1000,
		}
		for {
			response, err := client.ListTrips(cmd.Context(), &request)
			if err != nil {
				return err
			}
			for _, trip := range response.Trips {
				printJSON(cmd, trip)
				validate(cmd, trip)
			}
			if response.ContinuationToken == "" {
				break
			}
			request.ContinuationToken = response.ContinuationToken
		}
		return nil
	}
	return cmd
}

func newListFuelEventsCommand(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fuel-events [object-id]",
		Short:   "List fuel events for an object",
		GroupID: "fuel-events",
		Args:    cobra.ExactArgs(1),
	}
	fromTime := cmd.Flags().Time("from", time.Now().Add(-7*24*time.Hour), []string{time.RFC3339}, "From time (RFC3339 format)")
	toTime := cmd.Flags().Time("to", time.Now(), []string{time.RFC3339}, "To time (RFC3339 format)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd, cfg)
		if err != nil {
			return err
		}
		request := trusttrack.ListFuelEventsRequest{
			ObjectID: args[0],
			FromTime: *fromTime,
			ToTime:   *toTime,
			Limit:    1000,
		}
		for {
			response, err := client.ListFuelEvents(cmd.Context(), &request)
			if err != nil {
				return err
			}
			for _, fuelEvent := range response.FuelEvents {
				printJSON(cmd, fuelEvent)
				validate(cmd, fuelEvent)
			}
			if response.ContinuationToken == "" {
				break
			}
			request.ContinuationToken = response.ContinuationToken
		}
		return nil
	}
	return cmd
}

func promptSecret(cmd *cobra.Command, prompt string) (string, error) {
	cmd.Print(prompt)
	input, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	cmd.Println()
	return string(input), nil
}

func printJSON(_ *cobra.Command, msg proto.Message) {
	fmt.Println(protojson.Format(msg))
}

func validate(cmd *cobra.Command, msg proto.Message) {
	if v, _ := cmd.Flags().GetBool("validate"); v {
		style := lipgloss.NewStyle().Foreground(lipgloss.Red)
		if err := protovalidate.Validate(msg); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", style.Render(fmt.Sprintf("%v", err)))
		}
	}
}
