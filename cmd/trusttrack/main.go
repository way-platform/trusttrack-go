package main

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"time"

	"buf.build/go/protovalidate"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/trusttrack-go"
	"github.com/way-platform/trusttrack-go/cmd/trusttrack/internal/auth"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCommand(),
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			base := c(lipgloss.Black, lipgloss.White)
			baseInverted := c(lipgloss.White, lipgloss.Black)
			return fang.ColorScheme{
				Base:         base,
				Title:        base,
				Description:  base,
				Comment:      base,
				Flag:         base,
				FlagDefault:  base,
				Command:      base,
				QuotedString: base,
				Argument:     base,
				Help:         base,
				Dash:         base,
				ErrorHeader:  [2]color.Color{baseInverted, base},
				ErrorDetails: base,
			}
		}),
	); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trusttrack",
		Short: "TrustTrack API CLI",
	}
	cmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	cmd.PersistentFlags().Bool("validate", false, "Enable message validation")
	cmd.AddGroup(&cobra.Group{
		ID:    "objects",
		Title: "Objects",
	})
	cmd.AddCommand(newListObjectsCommand())
	cmd.AddCommand(newListObjectsLastPositionCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "object-groups",
		Title: "Object Groups",
	})
	cmd.AddCommand(newListObjectGroupsCommand())
	cmd.AddCommand(newGetObjectGroupCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "coordinates",
		Title: "Coordinates",
	})
	cmd.AddCommand(newListObjectCoordinatesCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "trips",
		Title: "Trips",
	})
	cmd.AddCommand(newListTripsCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "fuel-events",
		Title: "Fuel Events",
	})
	cmd.AddCommand(newListFuelEventsCommand())
	cmd.AddGroup(auth.NewGroup())
	cmd.AddCommand(auth.NewCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "utils",
		Title: "Utils",
	})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

func newClient(cmd *cobra.Command) (*trusttrack.Client, error) {
	debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
	client, err := auth.NewClient(
		trusttrack.WithDebug(debug),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newListObjectsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects",
		Short:   "List objects",
		GroupID: "objects",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
		if err != nil {
			return err
		}
		response, err := client.ListObjects(cmd.Context(), &trusttrack.ListObjectsRequest{})
		if err != nil {
			return err
		}
		for _, object := range response.Objects {
			fmt.Println(protojson.Format(object))
			validate(cmd, object)
		}
		return nil
	}
	return cmd
}

func newListObjectsLastPositionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects-last-position",
		Short:   "List objects with their last position",
		GroupID: "objects",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
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
				fmt.Println(protojson.Format(object))
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

func newListObjectCoordinatesCommand() *cobra.Command {
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
		client, err := newClient(cmd)
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
				fmt.Println(protojson.Format(coordinate))
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

func newListObjectGroupsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object-groups",
		Short:   "List object groups",
		GroupID: "object-groups",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
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
				fmt.Println(protojson.Format(objectGroup))
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

func newGetObjectGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object-group [external-id]",
		Short:   "Get a specific object group by external ID",
		GroupID: "object-groups",
		Args:    cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
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
		fmt.Println(protojson.Format(response.ObjectGroup))
		validate(cmd, response.ObjectGroup)
		return nil
	}
	return cmd
}

func newListTripsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trips [object-id]",
		Short:   "List trips for an object",
		GroupID: "trips",
		Args:    cobra.ExactArgs(1),
	}
	fromTime := cmd.Flags().Time("from", time.Now().Add(-24*time.Hour), []string{time.RFC3339}, "From time (RFC3339 format)")
	toTime := cmd.Flags().Time("to", time.Time{}, []string{time.RFC3339}, "To time (RFC3339 format, optional)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
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
				fmt.Println(protojson.Format(trip))
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

func newListFuelEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fuel-events [object-id]",
		Short:   "List fuel events for an object",
		GroupID: "fuel-events",
		Args:    cobra.ExactArgs(1),
	}
	fromTime := cmd.Flags().Time("from", time.Now().Add(-7*24*time.Hour), []string{time.RFC3339}, "From time (RFC3339 format)")
	toTime := cmd.Flags().Time("to", time.Now(), []string{time.RFC3339}, "To time (RFC3339 format)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
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
				fmt.Println(protojson.Format(fuelEvent))
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

func validate(cmd *cobra.Command, object proto.Message) {
	if validate, _ := cmd.Root().PersistentFlags().GetBool("validate"); validate {
		style := lipgloss.NewStyle().Foreground(lipgloss.Red)
		if err := protovalidate.Validate(object); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", style.Render(fmt.Sprintf("%v", err)))
		}
	}
}
