package cmd

import (
	"fmt"
	"os"

	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/spf13/cobra"
)

// settingCmd represents the setting command
var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Manage CLI settings",
	Long:  `Manage CLI settings including themes and other preferences.`,
}

// settingThemeCmd represents the setting theme command
var settingThemeCmd = &cobra.Command{
	Use:   "theme [theme-name]",
	Short: "View or set the theme",
	Long: `View the current theme or set a new theme. Themes customize the color 
scheme of the CLI interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			themeName := args[0]
			if err := config.SetTheme(themeName); err != nil {
				fmt.Fprintf(os.Stderr, "Error setting theme: %v\n", err)
				return
			}
			fmt.Printf("Theme set to: %s\n", themeName)
		} else {
			theme := config.GetTheme()
			fmt.Printf("Current theme: %s\n", theme)
		}
	},
}

func init() {
	rootCmd.AddCommand(settingCmd)
	settingCmd.AddCommand(settingThemeCmd)
}
