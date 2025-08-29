package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewCommand returns a new [cobra.Command] for CLI authentication.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Authenticate to the TrustTrack API",
		GroupID: "auth",
	}
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	return cmd
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the TrustTrack API",
	}
	apiKey := cmd.Flags().String("api-key", "-", "API key to use for authentication")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if *apiKey == "-" {
			cmd.Println("\nEnter API key:")
			input, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			*apiKey = string(input)
		}
		if err := writeFile(&File{
			APIKey: *apiKey,
		}); err != nil {
			return err
		}
		fmt.Println("\nLogged in to the TrustTrack API.")
		return nil
	}
	return cmd
}

func newLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from the TrustTrack API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := removeFile(); err != nil {
				return err
			}
			cmd.Println("Logged out.")
			return nil
		},
	}
}
