package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BitbucketAPIBaseURL = "https://api.bitbucket.org/2.0"
)

// BitbucketClient represents a Bitbucket API client
type BitbucketClient struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string // App password or token
}

// NewBitbucketClient creates a new Bitbucket API client
func NewBitbucketClient(username, password string) *BitbucketClient {
	return &BitbucketClient{
		baseURL: BitbucketAPIBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		username: username,
		password: password,
	}
}

func (c *BitbucketClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
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

	if c.username != "" && c.password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		req.Header.Set("Authorization", "Basic "+auth)
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
		return nil, fmt.Errorf("Bitbucket API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// BitbucketPullRequest represents a Bitbucket pull request
type BitbucketPullRequest struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	URL         string    `json:"links.html.href"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	Author      struct {
		DisplayName string `json:"display_name"`
		Username    string `json:"username"`
	} `json:"author"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
}

// BitbucketIssue represents a Bitbucket issue
type BitbucketIssue struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content.raw"`
	State       string    `json:"state"`
	Kind        string    `json:"kind"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	Reporter    struct {
		DisplayName string `json:"display_name"`
		Username    string `json:"username"`
	} `json:"reporter"`
}

// ListPullRequests lists pull requests for a repository
func (c *BitbucketClient) ListPullRequests(ctx context.Context, workspace, repo string, state string) ([]BitbucketPullRequest, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests", workspace, repo)
	if state != "" {
		path += "?state=" + state
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Values []BitbucketPullRequest `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Values, nil
}

// GetPullRequest gets a specific pull request
func (c *BitbucketClient) GetPullRequest(ctx context.Context, workspace, repo string, id int) (*BitbucketPullRequest, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", workspace, repo, id)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr BitbucketPullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// ListIssues lists issues for a repository
func (c *BitbucketClient) ListIssues(ctx context.Context, workspace, repo string, state string) ([]BitbucketIssue, error) {
	path := fmt.Sprintf("/repositories/%s/%s/issues", workspace, repo)
	if state != "" {
		path += "?state=" + state
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Values []BitbucketIssue `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Values, nil
}
