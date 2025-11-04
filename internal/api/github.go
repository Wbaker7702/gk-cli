package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	GitHubAPIBaseURL = "https://api.github.com"
)

// GitHubClient represents a GitHub API client
type GitHubClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		baseURL: GitHubAPIBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

func (c *GitHubClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID          int       `json:"id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	State       string    `json:"state"` // open, closed
	URL         string    `json:"html_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user"`
	Head        Branch    `json:"head"`
	Base        Branch    `json:"base"`
	Mergeable   *bool     `json:"mergeable"`
	MergeableState string `json:"mergeable_state"`
}

// User represents a GitHub user
type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Branch represents a GitHub branch
type Branch struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
	Repo Repo   `json:"repo"`
}

// Repo represents a GitHub repository
type Repo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
	Private     bool   `json:"private"`
}

// Issue represents a GitHub issue
type Issue struct {
	ID          int       `json:"id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	URL         string    `json:"html_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user"`
	Labels      []Label   `json:"labels"`
	PullRequest *struct {
		URL string `json:"url"`
	} `json:"pull_request,omitempty"`
}

// Label represents a GitHub label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ListPullRequests lists pull requests for a repository
func (c *GitHubClient) ListPullRequests(ctx context.Context, owner, repo string, state string) ([]PullRequest, error) {
	if state == "" {
		state = "open"
	}

	path := fmt.Sprintf("/repos/%s/%s/pulls?state=%s", owner, repo, state)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prs []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return prs, nil
}

// GetPullRequest gets a specific pull request
func (c *GitHubClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// ListIssues lists issues for a repository
func (c *GitHubClient) ListIssues(ctx context.Context, owner, repo string, state string) ([]Issue, error) {
	if state == "" {
		state = "open"
	}

	path := fmt.Sprintf("/repos/%s/%s/issues?state=%s", owner, repo, state)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issues, nil
}

// CreatePullRequestComment creates a comment on a pull request
func (c *GitHubClient) CreatePullRequestComment(ctx context.Context, owner, repo string, number int, body string) error {
	path := fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, number)
	payload := map[string]string{"body": body}
	resp, err := c.doRequest(ctx, "POST", path, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
