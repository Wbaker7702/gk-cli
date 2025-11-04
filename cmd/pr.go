package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gitkraken/gk-cli/internal/api"
	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/gitkraken/gk-cli/internal/workspace"
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	Long: `Manage pull requests across repositories in your workspace. Requires 
a Cloud Workspace with a connected provider.`,
}

// prListCmd represents the pr list command
var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pull requests",
	Long: `List all pull requests for all repositories in your workspace. 
Filters can be applied to show specific states or repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := getWorkspace()
		if err != nil {
			return err
		}

		state, _ := cmd.Flags().GetString("state")
		if state == "" {
			state = "open"
		}

		factory := api.NewProviderFactory()
		cfg := config.Get()

		// Setup providers from config
		if providers, ok := cfg.Providers["github"].(map[string]interface{}); ok {
			if token, ok := providers["token"].(string); ok {
				factory.SetGitHubToken(token)
			}
		}
		if providers, ok := cfg.Providers["gitlab"].(map[string]interface{}); ok {
			if token, ok := providers["token"].(string); ok {
				factory.SetGitLabToken(token)
			}
		}
		if providers, ok := cfg.Providers["bitbucket"].(map[string]interface{}); ok {
			if user, ok := providers["username"].(string); ok {
				if pass, ok := providers["password"].(string); ok {
					factory.SetBitbucketCreds(user, pass)
				}
			}
		}

		fmt.Printf("Listing pull requests (%s) for workspace '%s'...\n\n", state, ws.Name)

		totalPRs := 0
		for _, repo := range ws.Repos {
			if repo.Remote == "" {
				continue
			}

			providerName, owner, repoName, err := api.ParseRepoURL(repo.Remote)
			if err != nil {
				fmt.Printf("⚠ Skipping %s: %v\n", repo.Name, err)
				continue
			}

			provider, err := factory.GetProvider(providerName)
			if err != nil {
				fmt.Printf("⚠ Skipping %s: %v\n", repo.Name, err)
				continue
			}

			prs, err := provider.ListPullRequests(owner, repoName, state)
			if err != nil {
				fmt.Printf("⚠ Error fetching PRs for %s: %v\n", repo.Name, err)
				continue
			}

			if len(prs) > 0 {
				fmt.Printf("📦 %s/%s (%s):\n", owner, repoName, providerName)
				for _, pr := range prs {
					status := "🟢"
					if pr.State != "open" {
						status = "🔴"
					}
					fmt.Printf("  %s #%d: %s [%s → %s]\n", status, pr.Number, pr.Title, pr.SourceBranch, pr.TargetBranch)
					fmt.Printf("     Author: %s | %s\n", pr.Author, pr.URL)
				}
				fmt.Println()
				totalPRs += len(prs)
			}
		}

		if totalPRs == 0 {
			fmt.Println("No pull requests found.")
		} else {
			fmt.Printf("Total: %d pull request(s)\n", totalPRs)
		}

		return nil
	},
}

// prViewCmd represents the pr view command
var prViewCmd = &cobra.Command{
	Use:   "view [pr-number]",
	Short: "View pull request details",
	Long: `View detailed information about a specific pull request. 
If run in a git repository, it will detect the PR for that repo.
Otherwise, it will list all PRs for selection.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var prNumber int
		var err error

		if len(args) > 0 {
			prNumber, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
		}

		// Try to detect current repo
		cwd, _ := os.Getwd()
		repo, err := workspace.DetectRepo(cwd)
		if err != nil {
			return fmt.Errorf("not in a git repository. Please specify PR number or run from a git repo")
		}

		if repo.Remote == "" {
			return fmt.Errorf("repository has no remote URL")
		}

		providerName, owner, repoName, err := api.ParseRepoURL(repo.Remote)
		if err != nil {
			return fmt.Errorf("failed to parse repository URL: %w", err)
		}

		factory := api.NewProviderFactory()
		cfg := config.Get()

		// Setup provider
		switch providerName {
		case "github":
			if providers, ok := cfg.Providers["github"].(map[string]interface{}); ok {
				if token, ok := providers["token"].(string); ok {
					factory.SetGitHubToken(token)
				}
			}
		case "gitlab":
			if providers, ok := cfg.Providers["gitlab"].(map[string]interface{}); ok {
				if token, ok := providers["token"].(string); ok {
					factory.SetGitLabToken(token)
				}
			}
		case "bitbucket":
			if providers, ok := cfg.Providers["bitbucket"].(map[string]interface{}); ok {
				if user, ok := providers["username"].(string); ok {
					if pass, ok := providers["password"].(string); ok {
						factory.SetBitbucketCreds(user, pass)
					}
				}
			}
		}

		provider, err := factory.GetProvider(providerName)
		if err != nil {
			return fmt.Errorf("provider not configured: %w", err)
		}

		if prNumber == 0 {
			// List PRs for selection
			prs, err := provider.ListPullRequests(owner, repoName, "open")
			if err != nil {
				return fmt.Errorf("failed to list PRs: %w", err)
			}

			if len(prs) == 0 {
				fmt.Println("No open pull requests found.")
				return nil
			}

			fmt.Printf("Open pull requests for %s/%s:\n\n", owner, repoName)
			for i, pr := range prs {
				fmt.Printf("%d. #%d: %s\n", i+1, pr.Number, pr.Title)
			}
			fmt.Print("\nSelect PR number: ")
			var selection int
			fmt.Scanf("%d", &selection)
			if selection < 1 || selection > len(prs) {
				return fmt.Errorf("invalid selection")
			}
			prNumber = prs[selection-1].Number
		}

		pr, err := provider.GetPullRequest(owner, repoName, prNumber)
		if err != nil {
			return fmt.Errorf("failed to get PR: %w", err)
		}

		// Display PR details
		fmt.Printf("\n%s Pull Request #%d: %s\n", strings.ToUpper(providerName), pr.Number, pr.Title)
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("State:     %s\n", pr.State)
		fmt.Printf("Author:    %s\n", pr.Author)
		fmt.Printf("Branch:    %s → %s\n", pr.SourceBranch, pr.TargetBranch)
		fmt.Printf("URL:       %s\n", pr.URL)
		fmt.Printf("Created:   %s\n", pr.CreatedAt)
		fmt.Printf("Updated:   %s\n", pr.UpdatedAt)
		if pr.Body != "" {
			fmt.Println("\nDescription:")
			fmt.Println(strings.Repeat("-", 60))
			fmt.Println(pr.Body)
		}

		return nil
	},
}

// prSuggestCmd represents the pr suggest command
var prSuggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Create code suggestions for a pull request",
	Long: `Create code suggestions for an open pull request. This command should be 
run in a repository with an open pull request, after making local changes that 
you want to suggest.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PR suggest command - Creating code suggestions...")
		// TODO: Implement code suggestion creation
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(prListCmd)
	prCmd.AddCommand(prViewCmd)
	prCmd.AddCommand(prSuggestCmd)

	prListCmd.Flags().StringP("state", "s", "open", "Filter by state (open, closed, all)")
}
