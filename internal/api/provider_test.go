package api

import (
	"testing"
)

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		url      string
		provider string
		owner    string
		repo     string
		err      bool
	}{
		{
			url:      "https://github.com/user/repo.git",
			provider: "github",
			owner:    "user",
			repo:     "repo",
			err:      false,
		},
		{
			url:      "https://gitlab.com/namespace/project.git",
			provider: "gitlab",
			owner:    "namespace",
			repo:     "project",
			err:      false,
		},
		{
			url:      "https://bitbucket.org/workspace/repo.git",
			provider: "bitbucket",
			owner:    "workspace",
			repo:     "repo",
			err:      false,
		},
		{
			url:      "https://github.com/user/repo",
			provider: "github",
			owner:    "user",
			repo:     "repo",
			err:      false,
		},
		{
			url:      "invalid-url",
			provider: "",
			owner:    "",
			repo:     "",
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			provider, owner, repo, err := ParseRepoURL(tt.url)
			if tt.err && err == nil {
				t.Errorf("Expected error for URL %s", tt.url)
			}
			if !tt.err && err != nil {
				t.Errorf("Unexpected error for URL %s: %v", tt.url, err)
			}
			if provider != tt.provider {
				t.Errorf("Expected provider %s, got %s", tt.provider, provider)
			}
			if owner != tt.owner {
				t.Errorf("Expected owner %s, got %s", tt.owner, owner)
			}
			if repo != tt.repo {
				t.Errorf("Expected repo %s, got %s", tt.repo, repo)
			}
		})
	}
}
