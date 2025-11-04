package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/gitkraken/gk-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// providerCmd represents the provider command
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage hosting and issue providers",
	Long: `Manage connections to hosting and issue providers like GitHub, GitLab, 
and Bitbucket. These connections enable fetching pull requests and issues.`,
}

// providerAddCmd represents the provider add command
var providerAddCmd = &cobra.Command{
	Use:   "add [provider]",
	Short: "Add a provider connection",
	Long: `Add a provider connection (GitHub, GitLab, Bitbucket, etc.). 
For GitHub and GitLab, you can provide a personal access token.
For Bitbucket, provide username and app password.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName := strings.ToLower(args[0])
		cfg := config.Get()

		if cfg.Providers == nil {
			cfg.Providers = make(map[string]interface{})
		}

		switch providerName {
		case "github":
			token, _ := cmd.Flags().GetString("token")
			if token == "" {
				var err error
				token, err = utils.PromptString("GitHub Personal Access Token: ")
				if err != nil {
					return err
				}
			}
			cfg.Providers["github"] = map[string]interface{}{
				"token": token,
			}
			fmt.Println("✓ GitHub provider added")

		case "gitlab":
			token, _ := cmd.Flags().GetString("token")
			if token == "" {
				var err error
				token, err = utils.PromptString("GitLab Personal Access Token: ")
				if err != nil {
					return err
				}
			}
			cfg.Providers["gitlab"] = map[string]interface{}{
				"token": token,
			}
			fmt.Println("✓ GitLab provider added")

		case "bitbucket":
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			if username == "" {
				var err error
				username, err = utils.PromptString("Bitbucket Username: ")
				if err != nil {
					return err
				}
			}
			if password == "" {
				var err error
				password, err = utils.PromptString("Bitbucket App Password: ")
				if err != nil {
					return err
				}
			}
			cfg.Providers["bitbucket"] = map[string]interface{}{
				"username": username,
				"password": password,
			}
			fmt.Println("✓ Bitbucket provider added")

		default:
			return fmt.Errorf("unsupported provider: %s (supported: github, gitlab, bitbucket)", providerName)
		}

		// Update global config
		config.UpdateProviders(cfg.Providers)
		return nil
	},
}

// providerListCmd represents the provider list command
var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured providers",
	Long:  `List all configured hosting and issue providers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if cfg.Providers == nil || len(cfg.Providers) == 0 {
			fmt.Println("No providers configured.")
			fmt.Println("Add a provider with: gk provider add <github|gitlab|bitbucket>")
			return nil
		}

		fmt.Println("Configured providers:")
		for name, provider := range cfg.Providers {
			if providerMap, ok := provider.(map[string]interface{}); ok {
				fmt.Printf("  • %s", strings.Title(name))
				if name == "bitbucket" {
					if user, ok := providerMap["username"].(string); ok {
						fmt.Printf(" (user: %s)", user)
					}
				} else {
					if token, ok := providerMap["token"].(string); ok && len(token) > 0 {
						fmt.Printf(" (token: %s...)", token[:min(8, len(token))])
					}
				}
				fmt.Println()
			}
		}
		return nil
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// providerRemoveCmd represents the provider remove command
var providerRemoveCmd = &cobra.Command{
	Use:   "remove [provider]",
	Short: "Remove a provider connection",
	Long:  `Remove a provider connection.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName := strings.ToLower(args[0])
		cfg := config.Get()

		if cfg.Providers == nil {
			return fmt.Errorf("no providers configured")
		}

		if _, exists := cfg.Providers[providerName]; !exists {
			return fmt.Errorf("provider '%s' not found", providerName)
		}

		delete(cfg.Providers, providerName)
		config.UpdateProviders(cfg.Providers)

		fmt.Printf("✓ Removed provider: %s\n", providerName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)
	providerCmd.AddCommand(providerAddCmd)
	providerCmd.AddCommand(providerListCmd)
	providerCmd.AddCommand(providerRemoveCmd)

	providerAddCmd.Flags().StringP("token", "t", "", "Provider token (GitHub/GitLab)")
	providerAddCmd.Flags().StringP("username", "u", "", "Bitbucket username")
	providerAddCmd.Flags().StringP("password", "p", "", "Bitbucket app password")
}
