package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "ws",
	Short: "Manage workspaces",
	Long: `Manage GitKraken workspaces. Workspaces associate groups of repos and set 
the context for helpful commands that can operate on multiple repos at once.`,
}

// wsCreateCmd represents the ws create command
var wsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workspace",
	Long: `Create a new workspace. Workspaces can be local (existing only on your machine) 
or cloud (accessible on any machine and can be connected to hosting providers).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace create command")
		// TODO: Implement workspace creation
	},
}

// wsAddRepoCmd represents the ws add-repo command
var wsAddRepoCmd = &cobra.Command{
	Use:   "add-repo",
	Short: "Add a repository to a workspace",
	Long:  `Add a new repository to a workspace either by path or remote URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace add-repo command")
		// TODO: Implement adding repo to workspace
	},
}

// wsLocateCmd represents the ws locate command
var wsLocateCmd = &cobra.Command{
	Use:   "locate",
	Short: "Locate local repositories for a cloud workspace",
	Long: `If you're accessing a Cloud Workspace for the first time, you might need to 
locate the local repos on your machine. Run this command in the directory where 
your repos are located.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace locate command")
		// TODO: Implement repo location
	},
}

// wsCloneCmd represents the ws clone command
var wsCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone all repositories in a workspace",
	Long: `Clone all repos in a Workspace at once into a single directory. This is 
helpful for onboarding when your team works on multiple repos.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace clone command")
		// TODO: Implement workspace clone
	},
}

// wsInsightsCmd represents the ws insights command
var wsInsightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "View pull request insights for repositories in a workspace",
	Long: `See metrics for all repositories in a Cloud Workspace including:
- Average Cycle Time
- Average Throughput
- Merge Rate
- Opened/Merged/Closed Pull Requests`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace insights command")
		// TODO: Implement insights
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(wsCreateCmd)
	workspaceCmd.AddCommand(wsAddRepoCmd)
	workspaceCmd.AddCommand(wsLocateCmd)
	workspaceCmd.AddCommand(wsCloneCmd)
	workspaceCmd.AddCommand(wsInsightsCmd)

	// Git operations on workspaces
	gitOps := []string{"fetch", "pull", "push", "checkout"}
	for _, op := range gitOps {
		opCmd := &cobra.Command{
			Use:   op,
			Short: fmt.Sprintf("Perform '%s' operation on all repos in workspace", op),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("Workspace %s command\n", op)
				// TODO: Implement git operation on all repos
			},
		}
		workspaceCmd.AddCommand(opCmd)
	}
}
