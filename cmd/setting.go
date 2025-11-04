package cmd

import (
	"fmt"

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
			fmt.Printf("Setting theme to: %s\n", args[0])
			// TODO: Set theme in config
		} else {
			fmt.Println("Current theme: default")
			// TODO: Read theme from config
		}
	},
}

func init() {
	rootCmd.AddCommand(settingCmd)
	settingCmd.AddCommand(settingThemeCmd)
}
