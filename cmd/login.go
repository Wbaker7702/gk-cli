package cmd

import (
	"fmt"
	"os"

	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with GitKraken",
	Long: `Login to GitKraken to enable cloud features and workspace synchronization.
This will open a browser window for authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Login command - Opening browser for authentication...")
		fmt.Println("TODO: Implement OAuth flow")
		// TODO: Implement browser-based OAuth flow
		// 1. Start local server on port 1314
		// 2. Open browser to GitKraken OAuth URL
		// 3. Handle callback and store tokens
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
