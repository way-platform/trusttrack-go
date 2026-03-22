package main

import (
	"context"
	"image/color"
	"net/http"
	"os"

	"charm.land/fang/v2"
	"charm.land/lipgloss/v2"
	"github.com/adrg/xdg"
	trusttrack "github.com/way-platform/trusttrack-go"
	"github.com/way-platform/trusttrack-go/cli"
)

var debug bool

func main() {
	credPath, _ := xdg.ConfigFile("trusttrack-go/credentials.json")
	cmd := cli.NewCommand(
		cli.WithCredentialStore(cli.NewFileStore(credPath)),
		cli.WithHTTPClient(&http.Client{
			Transport: &trusttrack.DebugTransport{
				Enabled: &debug,
			},
		}),
	)
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	if err := fang.Execute(
		context.Background(),
		cmd,
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
