package cmd

import (
	"fmt"

	"github.com/gitkraken/gk-cli/internal/launchpad"
	"github.com/spf13/cobra"
)

// launchpadCmd represents the launchpad command
var launchpadCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "Open GitKraken Launchpad",
	Long: `GitKraken Launchpad is a unified dashboard that consolidates PRs, Issues, 
and WIPs across all of the repositories in a Cloud Workspace. You can view the 
details of any item and take action on your most important tasks.

Shortcuts:
  p - Pin/unpin items
  s - Snooze/unsnooze items`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := getWorkspace()
		if err != nil {
			return err
		}

		fmt.Println("Loading Launchpad...")
		lp, err := launchpad.LoadItems(ws)
		if err != nil {
			return fmt.Errorf("failed to load launchpad: %w", err)
		}

		lp.Display()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(launchpadCmd)
}
