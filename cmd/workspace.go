package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gitkraken/gk-cli/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	workspaceName string
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "ws",
	Short: "Manage workspaces",
	Long: `Manage GitKraken workspaces. Workspaces associate groups of repos and set 
the context for helpful commands that can operate on multiple repos at once.`,
}

// getWorkspace gets the workspace to use (from flag or prompts)
func getWorkspace() (*workspace.Workspace, error) {
	if workspaceName != "" {
		return workspace.Load(workspaceName)
	}

	// Try to get default workspace
	workspaces, err := workspace.List()
	if err != nil {
		return nil, err
	}

	if len(workspaces) == 0 {
		return nil, fmt.Errorf("no workspaces found. Create one with 'gk ws create'")
	}

	if len(workspaces) == 1 {
		return workspace.Load(workspaces[0])
	}

	// Multiple workspaces - prompt user
	fmt.Println("Available workspaces:")
	for i, ws := range workspaces {
		fmt.Printf("  %d. %s\n", i+1, ws)
	}
	fmt.Print("Select workspace (number or name): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Try as number first
	var selected string
	if idx := 0; len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
		if int(input[0]-'0') <= len(workspaces) {
			selected = workspaces[int(input[0]-'0')-1]
		}
	} else {
		selected = input
	}

	if selected == "" {
		return nil, fmt.Errorf("invalid workspace selection")
	}

	return workspace.Load(selected)
}

// wsCreateCmd represents the ws create command
var wsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new workspace",
	Long: `Create a new workspace. Workspaces can be local (existing only on your machine) 
or cloud (accessible on any machine and can be connected to hosting providers).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			fmt.Print("Workspace name: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
			if name == "" {
				return fmt.Errorf("workspace name cannot be empty")
			}
		}

		wsType, _ := cmd.Flags().GetString("type")
		if wsType == "" {
			wsType = "local" // Default to local
		}

		description, _ := cmd.Flags().GetString("description")

		ws, err := workspace.Create(name, wsType, description)
		if err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		fmt.Printf("✓ Created %s workspace '%s'\n", wsType, name)
		return nil
	},
}

// wsAddRepoCmd represents the ws add-repo command
var wsAddRepoCmd = &cobra.Command{
	Use:   "add-repo [path-or-url]",
	Short: "Add a repository to a workspace",
	Long:  `Add a new repository to a workspace either by path or remote URL.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := getWorkspace()
		if err != nil {
			return err
		}

		var repoPath string
		if len(args) > 0 {
			repoPath = args[0]
		} else {
			fmt.Print("Repository path or URL: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			repoPath = strings.TrimSpace(input)
		}

		if repoPath == "" {
			return fmt.Errorf("repository path cannot be empty")
		}

		// Check if it's a URL or local path
		var repo *workspace.Repo
		if strings.HasPrefix(repoPath, "http://") || strings.HasPrefix(repoPath, "https://") || strings.HasPrefix(repoPath, "git@") {
			// Remote URL - we'll need to clone it or just store the URL
			repo = &workspace.Repo{
				Remote: repoPath,
				Name:   filepath.Base(strings.TrimSuffix(repoPath, ".git")),
			}
		} else {
			// Local path - detect repo info
			absPath, err := filepath.Abs(repoPath)
			if err != nil {
				return fmt.Errorf("invalid path: %w", err)
			}
			repo, err = workspace.DetectRepo(absPath)
			if err != nil {
				return fmt.Errorf("failed to detect repository: %w", err)
			}
		}

		if err := ws.AddRepo(*repo); err != nil {
			return fmt.Errorf("failed to add repository: %w", err)
		}

		fmt.Printf("✓ Added repository '%s' to workspace '%s'\n", repo.Name, ws.Name)
		return nil
	},
}

// wsLocateCmd represents the ws locate command
var wsLocateCmd = &cobra.Command{
	Use:   "locate [directory]",
	Short: "Locate local repositories for a cloud workspace",
	Long: `If you're accessing a Cloud Workspace for the first time, you might need to 
locate the local repos on your machine. Run this command in the directory where 
your repos are located.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := getWorkspace()
		if err != nil {
			return err
		}

		var searchDir string
		if len(args) > 0 {
			searchDir = args[0]
		} else {
			var err error
			searchDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}

		absDir, err := filepath.Abs(searchDir)
		if err != nil {
			return fmt.Errorf("invalid directory: %w", err)
		}

		fmt.Printf("Searching for repositories in %s...\n", absDir)
		repos, err := workspace.LocateRepos(absDir)
		if err != nil {
			return fmt.Errorf("failed to locate repositories: %w", err)
		}

		if len(repos) == 0 {
			fmt.Println("No git repositories found")
			return nil
		}

		fmt.Printf("Found %d repositories:\n", len(repos))
		for _, repo := range repos {
			fmt.Printf("  - %s (%s)\n", repo.Name, repo.Path)
			ws.AddRepo(*repo)
		}

		fmt.Printf("✓ Added %d repositories to workspace '%s'\n", len(repos), ws.Name)
		return nil
	},
}

// wsCloneCmd represents the ws clone command
var wsCloneCmd = &cobra.Command{
	Use:   "clone [directory]",
	Short: "Clone all repositories in a workspace",
	Long: `Clone all repos in a Workspace at once into a single directory. This is 
helpful for onboarding when your team works on multiple repos.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := getWorkspace()
		if err != nil {
			return err
		}

		var targetDir string
		if len(args) > 0 {
			targetDir = args[0]
		} else {
			targetDir = ws.Name
		}

		absDir, err := filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("invalid directory: %w", err)
		}

		if err := os.MkdirAll(absDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		fmt.Printf("Cloning repositories to %s...\n", absDir)
		cloned := 0
		for _, repo := range ws.Repos {
			if repo.Remote == "" {
				fmt.Printf("⚠ Skipping %s (no remote URL)\n", repo.Name)
				continue
			}

			repoPath := filepath.Join(absDir, repo.Name)
			fmt.Printf("Cloning %s...\n", repo.Remote)
			if err := cloneRepo(repo.Remote, repoPath); err != nil {
				fmt.Printf("✗ Failed to clone %s: %v\n", repo.Name, err)
				continue
			}
			cloned++
		}

		fmt.Printf("✓ Cloned %d repositories\n", cloned)
		return nil
	},
}

func cloneRepo(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
		op := op // Capture loop variable
		opCmd := &cobra.Command{
			Use:   op + " [git-args...]",
			Short: fmt.Sprintf("Perform '%s' operation on all repos in workspace", op),
			RunE: func(cmd *cobra.Command, args []string) error {
				ws, err := getWorkspace()
				if err != nil {
					return err
				}

				fmt.Printf("Running 'git %s' on all repositories in workspace '%s'...\n", op, ws.Name)
				if err := workspace.GitOperation(ws, op, args); err != nil {
					return fmt.Errorf("git operation failed: %w", err)
				}

				fmt.Println("✓ Completed")
				return nil
			},
		}
		workspaceCmd.AddCommand(opCmd)
	}

	// Add flags
	workspaceCmd.PersistentFlags().StringVarP(&workspaceName, "workspace", "w", "", "workspace name")
	wsCreateCmd.Flags().StringP("type", "t", "local", "workspace type (local or cloud)")
	wsCreateCmd.Flags().StringP("description", "d", "", "workspace description")
}
