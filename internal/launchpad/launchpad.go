package launchpad

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gitkraken/gk-cli/internal/api"
	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/gitkraken/gk-cli/internal/workspace"
)

// Item represents a PR or Issue in the launchpad
type Item struct {
	Type        string // "pr" or "issue"
	Provider    string
	Repo        string
	Number      int
	Title       string
	State       string
	Author      string
	URL         string
	CreatedAt   string
	UpdatedAt   string
	Pinned      bool
	Snoozed     bool
}

// Launchpad represents the launchpad data
type Launchpad struct {
	Items []Item
}

// LoadItems loads PRs and Issues from workspace repositories
func LoadItems(ws *workspace.Workspace) (*Launchpad, error) {
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

	var items []Item
	ctx := context.Background()

	for _, repo := range ws.Repos {
		if repo.Remote == "" {
			continue
		}

		providerName, owner, repoName, err := api.ParseRepoURL(repo.Remote)
		if err != nil {
			continue
		}

		provider, err := factory.GetProvider(providerName)
		if err != nil {
			continue
		}

		// Get PRs
		prs, err := provider.ListPullRequests(owner, repoName, "open")
		if err == nil {
			for _, pr := range prs {
				items = append(items, Item{
					Type:      "pr",
					Provider:  providerName,
					Repo:      fmt.Sprintf("%s/%s", owner, repoName),
					Number:    pr.Number,
					Title:     pr.Title,
					State:     pr.State,
					Author:    pr.Author,
					URL:       pr.URL,
					CreatedAt: pr.CreatedAt,
					UpdatedAt: pr.UpdatedAt,
				})
			}
		}

		// Get Issues
		issues, err := provider.ListIssues(owner, repoName, "open")
		if err == nil {
			for _, issue := range issues {
				items = append(items, Item{
					Type:      "issue",
					Provider:  providerName,
					Repo:      fmt.Sprintf("%s/%s", owner, repoName),
					Number:    issue.Number,
					Title:     issue.Title,
					State:     issue.State,
					Author:    issue.Author,
					URL:       issue.URL,
					CreatedAt: issue.CreatedAt,
					UpdatedAt: issue.UpdatedAt,
				})
			}
		}
	}

	// Sort by updated time (most recent first)
	sort.Slice(items, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, items[i].UpdatedAt)
		timeJ, _ := time.Parse(time.RFC3339, items[j].UpdatedAt)
		return timeI.After(timeJ)
	})

	return &Launchpad{Items: items}, nil
}

// Display displays the launchpad items
func (l *Launchpad) Display() {
	if len(l.Items) == 0 {
		fmt.Println("No items found.")
		return
	}

	fmt.Println("🚀 GitKraken Launchpad")
	fmt.Println(strings.Repeat("=", 80))

	pinned := []Item{}
	regular := []Item{}

	for _, item := range l.Items {
		if item.Pinned {
			pinned = append(pinned, item)
		} else if !item.Snoozed {
			regular = append(regular, item)
		}
	}

	if len(pinned) > 0 {
		fmt.Println("\n📌 Pinned:")
		for i, item := range pinned {
			l.displayItem(i+1, item, true)
		}
	}

	if len(regular) > 0 {
		if len(pinned) > 0 {
			fmt.Println("\n📋 Open Items:")
		}
		for i, item := range regular {
			l.displayItem(i+1, item, false)
		}
	}

	fmt.Printf("\nTotal: %d items\n", len(l.Items))
	fmt.Println("\nCommands:")
	fmt.Println("  Press 'p' to pin/unpin an item")
	fmt.Println("  Press 's' to snooze/unsnooze an item")
	fmt.Println("  Press 'q' to quit")
}

func (l *Launchpad) displayItem(index int, item Item, pinned bool) {
	icon := "🔵"
	if item.Type == "issue" {
		icon = "⚪"
	}
	pinIcon := ""
	if pinned {
		pinIcon = "📌 "
	}

	fmt.Printf("%s%d. %s%s #%d: %s\n", pinIcon, index, icon, item.Number, item.Title)
	fmt.Printf("   %s/%s | %s | %s\n", item.Provider, item.Repo, item.Author, item.URL)
}
