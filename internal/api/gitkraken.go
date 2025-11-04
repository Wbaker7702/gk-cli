package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gitkraken/gk-cli/internal/auth"
)

const (
	// GitKraken API base URL (placeholder - actual URL needed)
	DefaultBaseURL = "https://api.gitkraken.com/v1"
)

// Client represents a GitKraken API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new GitKraken API client
func NewClient(baseURL string) (*Client, error) {
	token, err := auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}, nil
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
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

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, "PATCH", path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CloudPatch represents a GitKraken Cloud Patch
type CloudPatch struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Visibility  string    `json:"visibility"` // public, invite-only, private
}

// CreateCloudPatch creates a new cloud patch
func (c *Client) CreateCloudPatch(ctx context.Context, patchData []byte, name, description, visibility string) (*CloudPatch, error) {
	body := map[string]interface{}{
		"name":        name,
		"description": description,
		"visibility":  visibility,
		"patch_data":  string(patchData),
	}

	var result CloudPatch
	if err := c.Post(ctx, "/patches", body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCloudPatch retrieves a cloud patch by ID
func (c *Client) GetCloudPatch(ctx context.Context, patchID string) (*CloudPatch, error) {
	var result CloudPatch
	if err := c.Get(ctx, "/patches/"+patchID, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListCloudPatches lists all cloud patches
func (c *Client) ListCloudPatches(ctx context.Context) ([]CloudPatch, error) {
	var result struct {
		Patches []CloudPatch `json:"patches"`
	}
	if err := c.Get(ctx, "/patches", &result); err != nil {
		return nil, err
	}
	return result.Patches, nil
}

// DeleteCloudPatch deletes a cloud patch
func (c *Client) DeleteCloudPatch(ctx context.Context, patchID string) error {
	return c.Delete(ctx, "/patches/"+patchID)
}
