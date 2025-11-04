package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitOperation performs a git operation on all repos in a workspace
func GitOperation(ws *Workspace, operation string, args []string) error {
	repos := ws.GetRepoPaths()
	if len(repos) == 0 {
		return fmt.Errorf("no repositories in workspace '%s'", ws.Name)
	}

	var errors []string
	for _, repoPath := range repos {
		if err := runGitOperation(repoPath, operation, args); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", repoPath, err))
			continue
		}
		fmt.Printf("✓ %s\n", repoPath)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// runGitOperation runs a git command in a repository directory
func runGitOperation(repoPath, operation string, args []string) error {
	// Verify the path exists and is a git repository
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", repoPath)
	}

	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Build git command
	gitArgs := []string{operation}
	gitArgs = append(gitArgs, args...)

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DetectRepo detects repository information from a path
func DetectRepo(path string) (*Repo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if it's a git repository
	gitDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a git repository: %s", absPath)
	}

	// Get repository name from directory
	name := filepath.Base(absPath)

	// Try to get remote URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = absPath
	remoteBytes, err := cmd.Output()
	remote := strings.TrimSpace(string(remoteBytes))
	if err != nil {
		remote = "" // No remote is okay
	}

	repo := &Repo{
		Name:   name,
		Path:   absPath,
		Remote: remote,
	}

	// Detect provider from remote URL
	if remote != "" {
		if strings.Contains(remote, "github.com") {
			repo.Provider = "github"
		} else if strings.Contains(remote, "gitlab.com") {
			repo.Provider = "gitlab"
		} else if strings.Contains(remote, "bitbucket.org") {
			repo.Provider = "bitbucket"
		}
	}

	return repo, nil
}

// LocateRepos locates git repositories in a directory
func LocateRepos(rootDir string) ([]*Repo, error) {
	var repos []*Repo

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip hidden directories and .git directories
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") && info.Name() != ".") {
			return filepath.SkipDir
		}

		// Check if this directory is a git repository
		gitDir := filepath.Join(path, ".git")
		if info.IsDir() && info.Name() == ".git" {
			// Parent directory is a git repo
			repoPath := filepath.Dir(path)
			repo, err := DetectRepo(repoPath)
			if err == nil {
				repos = append(repos, repo)
			}
			return filepath.SkipDir // Don't walk into .git
		}

		return nil
	})

	return repos, err
}
