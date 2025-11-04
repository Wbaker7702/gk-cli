package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Workspace represents a GitKraken workspace
type Workspace struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` // "local" or "cloud"
	Repos       []Repo   `json:"repos"`
	Description string   `json:"description,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
}

// Repo represents a repository in a workspace
type Repo struct {
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`     // Local path
	Remote   string `json:"remote,omitempty"`    // Remote URL
	Provider string `json:"provider,omitempty"` // github, gitlab, bitbucket
}

var (
	workspacesDir string
)

// Init initializes the workspace system
func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	workspacesDir = filepath.Join(home, ".config", "gk", "workspaces")
	if err := os.MkdirAll(workspacesDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspaces directory: %w", err)
	}

	return nil
}

// Create creates a new workspace
func Create(name, workspaceType, description string) (*Workspace, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	if workspaceType != "local" && workspaceType != "cloud" {
		return nil, fmt.Errorf("workspace type must be 'local' or 'cloud'")
	}

	ws := &Workspace{
		Name:        name,
		Type:        workspaceType,
		Repos:       []Repo{},
		Description: description,
	}

	if err := ws.Save(); err != nil {
		return nil, err
	}

	return ws, nil
}

// Load loads a workspace by name
func Load(name string) (*Workspace, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	path := filepath.Join(workspacesDir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("workspace '%s' not found: %w", name, err)
	}

	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("failed to parse workspace file: %w", err)
	}

	return &ws, nil
}

// List lists all workspaces
func List() ([]string, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(workspacesDir)
	if err != nil {
		return []string{}, nil // No workspaces yet
	}

	var workspaces []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			name := strings.TrimSuffix(entry.Name(), ".json")
			workspaces = append(workspaces, name)
		}
	}

	return workspaces, nil
}

// Save saves the workspace to disk
func (ws *Workspace) Save() error {
	if err := Init(); err != nil {
		return err
	}

	path := filepath.Join(workspacesDir, ws.Name+".json")
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workspace: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write workspace file: %w", err)
	}

	return nil
}

// Delete deletes a workspace
func Delete(name string) error {
	if err := Init(); err != nil {
		return err
	}

	path := filepath.Join(workspacesDir, name+".json")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	return nil
}

// AddRepo adds a repository to the workspace
func (ws *Workspace) AddRepo(repo Repo) error {
	// Check if repo already exists
	for i, r := range ws.Repos {
		if r.Name == repo.Name || r.Path == repo.Path || r.Remote == repo.Remote {
			ws.Repos[i] = repo // Update existing
			return ws.Save()
		}
	}

	ws.Repos = append(ws.Repos, repo)
	return ws.Save()
}

// RemoveRepo removes a repository from the workspace
func (ws *Workspace) RemoveRepo(name string) error {
	for i, r := range ws.Repos {
		if r.Name == name {
			ws.Repos = append(ws.Repos[:i], ws.Repos[i+1:]...)
			return ws.Save()
		}
	}
	return fmt.Errorf("repository '%s' not found in workspace", name)
}

// GetRepoPaths returns all repository paths in the workspace
func (ws *Workspace) GetRepoPaths() []string {
	var paths []string
	for _, repo := range ws.Repos {
		if repo.Path != "" {
			paths = append(paths, repo.Path)
		}
	}
	return paths
}
