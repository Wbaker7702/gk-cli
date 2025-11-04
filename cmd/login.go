package cmd

import (
	"fmt"
	"os"

	"github.com/gitkraken/gk-cli/internal/auth"
	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with GitKraken",
	Long: `Login to GitKraken to enable cloud features and workspace synchronization.
This will open a browser window for authentication.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if already authenticated
		if config.IsAuthenticated() {
			fmt.Println("Already authenticated. Use 'gk logout' to logout first.")
			return nil
		}

		// TODO: Get actual client ID and secret from environment or config
		// For now, these need to be provided by GitKraken
		clientID := os.Getenv("GITKRAKEN_CLIENT_ID")
		clientSecret := os.Getenv("GITKRAKEN_CLIENT_SECRET")
		
		if clientID == "" || clientSecret == "" {
			fmt.Println("⚠ OAuth credentials not configured.")
			fmt.Println("Set GITKRAKEN_CLIENT_ID and GITKRAKEN_CLIENT_SECRET environment variables.")
			fmt.Println("\nFor now, authentication is not fully implemented.")
			fmt.Println("Once GitKraken OAuth credentials are available, this will:")
			fmt.Println("  1. Open your browser for authentication")
			fmt.Println("  2. Handle the callback on port 1314")
			fmt.Println("  3. Store your authentication tokens")
			return nil
		}

		// Initialize OAuth
		auth.InitOAuth(clientID, clientSecret, "http://localhost:1314/callback")

		fmt.Println("Opening browser for authentication...")
		if err := auth.StartAuthFlow(); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		return nil
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GitKraken",
	Long:  `Logout from GitKraken and clear stored authentication tokens.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.ClearAuth(); err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing authentication: %v\n", err)
			return
		}
		fmt.Println("Successfully logged out")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}
