package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gitkraken/gk-cli/internal/api"
	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/gitkraken/gk-cli/internal/patch"
	"github.com/gitkraken/gk-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Manage Cloud Patches",
	Long: `Manage Cloud Patches - Git patches that GitKraken securely stores for you 
so they can be easily shared with others across GitKraken CLI, GitKraken Desktop, 
and GitLens.`,
}

// patchCreateCmd represents the patch create command
var patchCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Cloud Patch",
	Long: `Create a Cloud Patch from your current changes. You will be prompted to 
provide information about the patch and sharing options (Public, Invite Only, Private).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.IsAuthenticated() {
			return fmt.Errorf("authentication required. Run 'gk login' first")
		}

		// Get patch data
		patchData, err := patch.CreatePatch("", "")
		if err != nil {
			return fmt.Errorf("failed to create patch: %w", err)
		}

		// Get patch name
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name, err = utils.PromptString("Patch name: ")
			if err != nil {
				return err
			}
		}

		// Get description
		description, _ := cmd.Flags().GetString("description")
		if description == "" {
			description, err = utils.PromptString("Description (optional): ")
			if err != nil {
				description = ""
			}
		}

		// Get visibility
		visibility, _ := cmd.Flags().GetString("visibility")
		if visibility == "" {
			choices := []string{"public", "invite-only", "private"}
			idx, err := utils.PromptChoice("Visibility:", choices)
			if err != nil {
				return err
			}
			visibility = choices[idx]
		}

		// Create cloud patch via API
		client, err := api.NewClient("")
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		ctx := context.Background()
		cloudPatch, err := client.CreateCloudPatch(ctx, patchData, name, description, visibility)
		if err != nil {
			return fmt.Errorf("failed to create cloud patch: %w", err)
		}

		fmt.Printf("✓ Cloud Patch created successfully!\n")
		fmt.Printf("  Name: %s\n", cloudPatch.Name)
		fmt.Printf("  URL: %s\n", cloudPatch.URL)
		fmt.Printf("  Visibility: %s\n", cloudPatch.Visibility)

		return nil
	},
}

// patchApplyCmd represents the patch apply command
var patchApplyCmd = &cobra.Command{
	Use:   "apply [patch-url]",
	Short: "Apply a Cloud Patch",
	Long: `Apply a Cloud Patch to the current repository. Can be applied to the 
working tree or to a new or existing branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf("Applying patch from: %s\n", args[0])
		} else {
			fmt.Println("Patch apply command - Please provide patch URL")
		}
		// TODO: Implement patch application
	},
}

// patchViewCmd represents the patch view command
var patchViewCmd = &cobra.Command{
	Use:   "view [patch-id]",
	Short: "Preview a Cloud Patch",
	Long:  `Preview the changes of a Cloud Patch without applying it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.IsAuthenticated() {
			return fmt.Errorf("authentication required. Run 'gk login' first")
		}

		patchID := args[0]
		client, err := api.NewClient("")
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		ctx := context.Background()
		cloudPatch, err := client.GetCloudPatch(ctx, patchID)
		if err != nil {
			return fmt.Errorf("failed to get patch: %w", err)
		}

		fmt.Printf("Cloud Patch: %s\n", cloudPatch.Name)
		if cloudPatch.Description != "" {
			fmt.Printf("Description: %s\n", cloudPatch.Description)
		}
		fmt.Printf("URL: %s\n", cloudPatch.URL)
		fmt.Printf("Visibility: %s\n", cloudPatch.Visibility)
		fmt.Println("\nNote: Patch content preview requires API support")

		return nil
	},
}

// patchListCmd represents the patch list command
var patchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Cloud Patches",
	Long:  `List all Cloud Patches you have created or have access to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.IsAuthenticated() {
			return fmt.Errorf("authentication required. Run 'gk login' first")
		}

		client, err := api.NewClient("")
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		ctx := context.Background()
		patches, err := client.ListCloudPatches(ctx)
		if err != nil {
			return fmt.Errorf("failed to list patches: %w", err)
		}

		if len(patches) == 0 {
			fmt.Println("No cloud patches found.")
			return nil
		}

		fmt.Println("Cloud Patches:")
		for _, p := range patches {
			fmt.Printf("  • %s (%s)\n", p.Name, p.Visibility)
			if p.Description != "" {
				fmt.Printf("    %s\n", p.Description)
			}
			fmt.Printf("    URL: %s\n", p.URL)
			fmt.Println()
		}

		return nil
	},
}

// patchDeleteCmd represents the patch delete command
var patchDeleteCmd = &cobra.Command{
	Use:   "delete [patch-id]",
	Short: "Delete a Cloud Patch",
	Long:  `Delete a Cloud Patch.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !config.IsAuthenticated() {
			return fmt.Errorf("authentication required. Run 'gk login' first")
		}

		patchID := args[0]
		confirmed, _ := utils.PromptYesNo(fmt.Sprintf("Delete patch %s?", patchID), false)
		if !confirmed {
			fmt.Println("Cancelled")
			return nil
		}

		client, err := api.NewClient("")
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		ctx := context.Background()
		if err := client.DeleteCloudPatch(ctx, patchID); err != nil {
			return fmt.Errorf("failed to delete patch: %w", err)
		}

		fmt.Printf("✓ Deleted patch %s\n", patchID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)
	patchCmd.AddCommand(patchCreateCmd)
	patchCmd.AddCommand(patchApplyCmd)
	patchCmd.AddCommand(patchViewCmd)
	patchCmd.AddCommand(patchListCmd)
	patchCmd.AddCommand(patchDeleteCmd)

	patchCreateCmd.Flags().StringP("name", "n", "", "Patch name")
	patchCreateCmd.Flags().StringP("description", "d", "", "Patch description")
	patchCreateCmd.Flags().StringP("visibility", "v", "", "Visibility (public, invite-only, private)")
}
