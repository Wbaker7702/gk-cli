package cmd

import (
	"fmt"
	"os"

	"github.com/gitkraken/gk-cli/internal/workspace"
	"github.com/spf13/cobra"
)

// wsListCmd represents the ws list command
var wsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Long:  `List all available workspaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaces, err := workspace.List()
		if err != nil {
			return fmt.Errorf("failed to list workspaces: %w", err)
		}

		if len(workspaces) == 0 {
			fmt.Println("No workspaces found. Create one with 'gk ws create'")
			return nil
		}

		fmt.Println("Workspaces:")
		for _, name := range workspaces {
			ws, err := workspace.Load(name)
			if err != nil {
				fmt.Printf("  - %s (error loading)\n", name)
				continue
			}
			fmt.Printf("  - %s (%s, %d repos)\n", name, ws.Type, len(ws.Repos))
		}

		return nil
	},
}

// wsShowCmd represents the ws show command
var wsShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show workspace details",
	Long:  `Show detailed information about a workspace.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			name = workspaceName
		}

		if name == "" {
			return fmt.Errorf("workspace name required")
		}

		ws, err := workspace.Load(name)
		if err != nil {
			return fmt.Errorf("failed to load workspace: %w", err)
		}

		fmt.Printf("Workspace: %s\n", ws.Name)
		fmt.Printf("Type: %s\n", ws.Type)
		if ws.Description != "" {
			fmt.Printf("Description: %s\n", ws.Description)
		}
		fmt.Printf("Repositories: %d\n", len(ws.Repos))
		if len(ws.Repos) > 0 {
			fmt.Println("\nRepositories:")
			for _, repo := range ws.Repos {
				if repo.Path != "" {
					fmt.Printf("  - %s (%s)\n", repo.Name, repo.Path)
				} else if repo.Remote != "" {
					fmt.Printf("  - %s (%s)\n", repo.Name, repo.Remote)
				} else {
					fmt.Printf("  - %s\n", repo.Name)
				}
			}
		}

		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(wsListCmd)
	workspaceCmd.AddCommand(wsShowCmd)
}
