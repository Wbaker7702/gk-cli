package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateWorkspace(t *testing.T) {
	// Setup test directory
	testDir := filepath.Join(os.TempDir(), "gk-test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Override workspaces directory
	workspacesDir = filepath.Join(testDir, "workspaces")

	ws, err := Create("test-workspace", "local", "Test workspace")
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	if ws.Name != "test-workspace" {
		t.Errorf("Expected name 'test-workspace', got '%s'", ws.Name)
	}

	if ws.Type != "local" {
		t.Errorf("Expected type 'local', got '%s'", ws.Type)
	}

	// Verify workspace can be loaded
	loaded, err := Load("test-workspace")
	if err != nil {
		t.Fatalf("Failed to load workspace: %v", err)
	}

	if loaded.Name != ws.Name {
		t.Errorf("Loaded workspace name mismatch")
	}
}

func TestAddRepo(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "gk-test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	workspacesDir = filepath.Join(testDir, "workspaces")

	ws, _ := Create("test-workspace", "local", "")
	
	repo := Repo{
		Name:   "test-repo",
		Path:   "/path/to/repo",
		Remote: "https://github.com/user/repo.git",
	}

	err := ws.AddRepo(repo)
	if err != nil {
		t.Fatalf("Failed to add repo: %v", err)
	}

	if len(ws.Repos) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(ws.Repos))
	}

	if ws.Repos[0].Name != "test-repo" {
		t.Errorf("Repo name mismatch")
	}
}

func TestListWorkspaces(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "gk-test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	workspacesDir = filepath.Join(testDir, "workspaces")

	Create("ws1", "local", "")
	Create("ws2", "local", "")

	workspaces, err := List()
	if err != nil {
		t.Fatalf("Failed to list workspaces: %v", err)
	}

	if len(workspaces) != 2 {
		t.Errorf("Expected 2 workspaces, got %d", len(workspaces))
	}
}
