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
	GitLabAPIBaseURL = "https://gitlab.com/api/v4"
)

// GitLabClient represents a GitLab API client
type GitLabClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewGitLabClient creates a new GitLab API client
func NewGitLabClient(token string) *GitLabClient {
	return &GitLabClient{
		baseURL: GitLabAPIBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

func (c *GitLabClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
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
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
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
		return nil, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// GitLabMergeRequest represents a GitLab merge request
type GitLabMergeRequest struct {
	ID          int       `json:"id"`
	IID         int       `json:"iid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	URL         string    `json:"web_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      GitLabUser `json:"author"`
	SourceBranch string   `json:"source_branch"`
	TargetBranch string   `json:"target_branch"`
}

// GitLabUser represents a GitLab user
type GitLabUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// GitLabIssue represents a GitLab issue
type GitLabIssue struct {
	ID          int       `json:"id"`
	IID         int       `json:"iid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	URL         string    `json:"web_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      GitLabUser `json:"author"`
	Labels      []string   `json:"labels"`
}

// ListMergeRequests lists merge requests for a project
func (c *GitLabClient) ListMergeRequests(ctx context.Context, projectID string, state string) ([]GitLabMergeRequest, error) {
	if state == "" {
		state = "opened"
	}

	path := fmt.Sprintf("/projects/%s/merge_requests?state=%s", projectID, state)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mrs []GitLabMergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mrs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return mrs, nil
}

// GetMergeRequest gets a specific merge request
func (c *GitLabClient) GetMergeRequest(ctx context.Context, projectID string, iid int) (*GitLabMergeRequest, error) {
	path := fmt.Sprintf("/projects/%s/merge_requests/%d", projectID, iid)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mr GitLabMergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &mr, nil
}

// ListIssues lists issues for a project
func (c *GitLabClient) ListIssues(ctx context.Context, projectID string, state string) ([]GitLabIssue, error) {
	if state == "" {
		state = "opened"
	}

	path := fmt.Sprintf("/projects/%s/issues?state=%s", projectID, state)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues []GitLabIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issues, nil
}
