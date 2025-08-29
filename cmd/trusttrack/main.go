package main

import (
	"context"
	"fmt"
	"image/color"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/trusttrack-go"
	"github.com/way-platform/trusttrack-go/cmd/trusttrack/internal/auth"
	"google.golang.org/protobuf/encoding/protojson"
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
	cmd.AddGroup(&cobra.Group{
		ID:    "objects",
		Title: "Objects",
	})
	cmd.AddCommand(newListObjectsCommand())
	cmd.AddCommand(newListObjectsLastCoordinateCommand())
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

func newListObjectsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects",
		Short:   "List objects",
		GroupID: "objects",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.ListObjects(context.Background(), &trusttrack.ListObjectsRequest{})
		if err != nil {
			return err
		}
		for _, object := range response.Objects {
			fmt.Println(protojson.Format(object))
		}
		return nil
	}
	return cmd
}

func newListObjectsLastCoordinateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects-last-coordinate",
		Short:   "List objects with their last coordinate",
		GroupID: "objects",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		request := trusttrack.ListObjectsLastCoordinateRequest{
			Limit: 1000,
		}
		for {
			response, err := client.ListObjectsLastCoordinate(context.Background(), &request)
			if err != nil {
				return err
			}
			for _, object := range response.Objects {
				fmt.Println(protojson.Format(object))
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
