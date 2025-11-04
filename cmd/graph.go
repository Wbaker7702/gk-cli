package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Open visual commit graph",
	Long: `Open a visual graph of the repository in your current directory in either 
GitKraken Desktop or GitLens in VS Code.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Opening visual commit graph...")
		// TODO: Implement graph opening
		// 1. Detect if GitKraken Desktop is installed
		// 2. Detect if GitLens is available
		// 3. Open appropriate application with repository path
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
