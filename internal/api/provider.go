package api

import (
	"fmt"
	"strings"
)

// Provider represents a git hosting provider
type Provider interface {
	GetName() string
	ListPullRequests(owner, repo, state string) ([]PullRequest, error)
	GetPullRequest(owner, repo string, number int) (*PullRequest, error)
	ListIssues(owner, repo, state string) ([]Issue, error)
}

// PullRequest is a unified pull request structure
type PullRequest struct {
	Provider    string
	ID          string
	Number      int
	Title       string
	Body        string
	State       string
	URL         string
	Author      string
	SourceBranch string
	TargetBranch string
	CreatedAt   string
	UpdatedAt   string
}

// Issue is a unified issue structure
type Issue struct {
	Provider string
	ID       string
	Number   int
	Title    string
	Body     string
	State    string
	URL      string
	Author   string
	Labels   []string
	CreatedAt string
	UpdatedAt string
}

// ProviderFactory creates provider clients
type ProviderFactory struct {
	githubToken   string
	gitlabToken   string
	bitbucketUser string
	bitbucketPass string
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

// SetGitHubToken sets the GitHub token
func (f *ProviderFactory) SetGitHubToken(token string) {
	f.githubToken = token
}

// SetGitLabToken sets the GitLab token
func (f *ProviderFactory) SetGitLabToken(token string) {
	f.gitlabToken = token
}

// SetBitbucketCreds sets Bitbucket credentials
func (f *ProviderFactory) SetBitbucketCreds(username, password string) {
	f.bitbucketUser = username
	f.bitbucketPass = password
}

// GetProvider gets a provider client by name
func (f *ProviderFactory) GetProvider(name string) (Provider, error) {
	switch strings.ToLower(name) {
	case "github":
		if f.githubToken == "" {
			return nil, fmt.Errorf("GitHub token not configured")
		}
		return &GitHubProviderAdapter{client: NewGitHubClient(f.githubToken)}, nil
	case "gitlab":
		if f.gitlabToken == "" {
			return nil, fmt.Errorf("GitLab token not configured")
		}
		return &GitLabProviderAdapter{client: NewGitLabClient(f.gitlabToken)}, nil
	case "bitbucket":
		if f.bitbucketUser == "" || f.bitbucketPass == "" {
			return nil, fmt.Errorf("Bitbucket credentials not configured")
		}
		return &BitbucketProviderAdapter{client: NewBitbucketClient(f.bitbucketUser, f.bitbucketPass)}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

// ParseRepoURL parses a repository URL to extract provider, owner, and repo
func ParseRepoURL(url string) (provider, owner, repo string, err error) {
	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	if strings.Contains(url, "github.com") {
		parts := strings.Split(url, "github.com/")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid GitHub URL")
		}
		repoParts := strings.Split(parts[1], "/")
		if len(repoParts) < 2 {
			return "", "", "", fmt.Errorf("invalid GitHub URL format")
		}
		return "github", repoParts[0], repoParts[1], nil
	}

	if strings.Contains(url, "gitlab.com") {
		parts := strings.Split(url, "gitlab.com/")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid GitLab URL")
		}
		repoParts := strings.Split(parts[1], "/")
		if len(repoParts) < 2 {
			return "", "", "", fmt.Errorf("invalid GitLab URL format")
		}
		// GitLab uses project ID which can be namespace/project
		return "gitlab", strings.Join(repoParts[:len(repoParts)-1], "/"), repoParts[len(repoParts)-1], nil
	}

	if strings.Contains(url, "bitbucket.org") {
		parts := strings.Split(url, "bitbucket.org/")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid Bitbucket URL")
		}
		repoParts := strings.Split(parts[1], "/")
		if len(repoParts) < 2 {
			return "", "", "", fmt.Errorf("invalid Bitbucket URL format")
		}
		return "bitbucket", repoParts[0], repoParts[1], nil
	}

	return "", "", "", fmt.Errorf("unsupported provider URL: %s", url)
}
