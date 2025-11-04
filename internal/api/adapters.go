package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GitHubProviderAdapter adapts GitHub client to Provider interface
type GitHubProviderAdapter struct {
	client *GitHubClient
}

func (a *GitHubProviderAdapter) GetName() string {
	return "github"
}

func (a *GitHubProviderAdapter) ListPullRequests(owner, repo, state string) ([]PullRequest, error) {
	ctx := context.Background()
	prs, err := a.client.ListPullRequests(ctx, owner, repo, state)
	if err != nil {
		return nil, err
	}

	result := make([]PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = PullRequest{
			Provider:      "github",
			ID:            strconv.Itoa(pr.ID),
			Number:        pr.Number,
			Title:         pr.Title,
			Body:          pr.Body,
			State:         pr.State,
			URL:           pr.URL,
			Author:        pr.User.Login,
			SourceBranch:  pr.Head.Ref,
			TargetBranch:  pr.Base.Ref,
			CreatedAt:     pr.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     pr.UpdatedAt.Format(time.RFC3339),
		}
	}
	return result, nil
}

func (a *GitHubProviderAdapter) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	ctx := context.Background()
	pr, err := a.client.GetPullRequest(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	return &PullRequest{
		Provider:      "github",
		ID:            strconv.Itoa(pr.ID),
		Number:        pr.Number,
		Title:         pr.Title,
		Body:          pr.Body,
		State:         pr.State,
		URL:           pr.URL,
		Author:        pr.User.Login,
		SourceBranch:  pr.Head.Ref,
		TargetBranch:  pr.Base.Ref,
		CreatedAt:     pr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     pr.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (a *GitHubProviderAdapter) ListIssues(owner, repo, state string) ([]Issue, error) {
	ctx := context.Background()
	issues, err := a.client.ListIssues(ctx, owner, repo, state)
	if err != nil {
		return nil, err
	}

	result := make([]Issue, 0)
	for _, issue := range issues {
		// Skip PRs (they have PullRequest field)
		if issue.PullRequest != nil {
			continue
		}

		labels := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			labels[i] = label.Name
		}

		result = append(result, Issue{
			Provider:   "github",
			ID:         strconv.Itoa(issue.ID),
			Number:     issue.Number,
			Title:      issue.Title,
			Body:       issue.Body,
			State:      issue.State,
			URL:        issue.URL,
			Author:     issue.User.Login,
			Labels:     labels,
			CreatedAt:  issue.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  issue.UpdatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

// GitLabProviderAdapter adapts GitLab client to Provider interface
type GitLabProviderAdapter struct {
	client *GitLabClient
}

func (a *GitLabProviderAdapter) GetName() string {
	return "gitlab"
}

func (a *GitLabProviderAdapter) ListPullRequests(owner, repo, state string) ([]PullRequest, error) {
	ctx := context.Background()
	projectID := owner + "/" + repo
	mrs, err := a.client.ListMergeRequests(ctx, projectID, state)
	if err != nil {
		return nil, err
	}

	result := make([]PullRequest, len(mrs))
	for i, mr := range mrs {
		result[i] = PullRequest{
			Provider:      "gitlab",
			ID:            strconv.Itoa(mr.ID),
			Number:        mr.IID,
			Title:         mr.Title,
			Body:          mr.Description,
			State:         mr.State,
			URL:           mr.URL,
			Author:        mr.Author.Username,
			SourceBranch:  mr.SourceBranch,
			TargetBranch:  mr.TargetBranch,
			CreatedAt:     mr.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     mr.UpdatedAt.Format(time.RFC3339),
		}
	}
	return result, nil
}

func (a *GitLabProviderAdapter) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	ctx := context.Background()
	projectID := owner + "/" + repo
	mr, err := a.client.GetMergeRequest(ctx, projectID, number)
	if err != nil {
		return nil, err
	}

	return &PullRequest{
		Provider:      "gitlab",
		ID:            strconv.Itoa(mr.ID),
		Number:        mr.IID,
		Title:         mr.Title,
		Body:          mr.Description,
		State:         mr.State,
		URL:           mr.URL,
		Author:        mr.Author.Username,
		SourceBranch:  mr.SourceBranch,
		TargetBranch:  mr.TargetBranch,
		CreatedAt:     mr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     mr.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (a *GitLabProviderAdapter) ListIssues(owner, repo, state string) ([]Issue, error) {
	ctx := context.Background()
	projectID := owner + "/" + repo
	issues, err := a.client.ListIssues(ctx, projectID, state)
	if err != nil {
		return nil, err
	}

	result := make([]Issue, len(issues))
	for i, issue := range issues {
		result[i] = Issue{
			Provider:   "gitlab",
			ID:         strconv.Itoa(issue.ID),
			Number:     issue.IID,
			Title:      issue.Title,
			Body:       issue.Description,
			State:      issue.State,
			URL:        issue.URL,
			Author:     issue.Author.Username,
			Labels:     issue.Labels,
			CreatedAt:  issue.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  issue.UpdatedAt.Format(time.RFC3339),
		}
	}
	return result, nil
}

// BitbucketProviderAdapter adapts Bitbucket client to Provider interface
type BitbucketProviderAdapter struct {
	client *BitbucketClient
}

func (a *BitbucketProviderAdapter) GetName() string {
	return "bitbucket"
}

func (a *BitbucketProviderAdapter) ListPullRequests(owner, repo, state string) ([]PullRequest, error) {
	ctx := context.Background()
	prs, err := a.client.ListPullRequests(ctx, owner, repo, state)
	if err != nil {
		return nil, err
	}

	result := make([]PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = PullRequest{
			Provider:      "bitbucket",
			ID:            strconv.Itoa(pr.ID),
			Number:        pr.ID,
			Title:         pr.Title,
			Body:          pr.Description,
			State:         strings.ToLower(pr.State),
			URL:           pr.URL,
			Author:        pr.Author.Username,
			SourceBranch:  pr.Source.Branch.Name,
			TargetBranch:  pr.Destination.Branch.Name,
			CreatedAt:     pr.CreatedOn.Format(time.RFC3339),
			UpdatedAt:     pr.UpdatedOn.Format(time.RFC3339),
		}
	}
	return result, nil
}

func (a *BitbucketProviderAdapter) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	ctx := context.Background()
	pr, err := a.client.GetPullRequest(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	return &PullRequest{
		Provider:      "bitbucket",
		ID:            strconv.Itoa(pr.ID),
		Number:        pr.ID,
		Title:         pr.Title,
		Body:          pr.Description,
		State:         strings.ToLower(pr.State),
		URL:           pr.URL,
		Author:        pr.Author.Username,
		SourceBranch:  pr.Source.Branch.Name,
		TargetBranch:  pr.Destination.Branch.Name,
		CreatedAt:     pr.CreatedOn.Format(time.RFC3339),
		UpdatedAt:     pr.UpdatedOn.Format(time.RFC3339),
	}, nil
}

func (a *BitbucketProviderAdapter) ListIssues(owner, repo, state string) ([]Issue, error) {
	ctx := context.Background()
	issues, err := a.client.ListIssues(ctx, owner, repo, state)
	if err != nil {
		return nil, err
	}

	result := make([]Issue, len(issues))
	for i, issue := range issues {
		result[i] = Issue{
			Provider:   "bitbucket",
			ID:         strconv.Itoa(issue.ID),
			Number:     issue.ID,
			Title:      issue.Title,
			Body:       issue.Content,
			State:      strings.ToLower(issue.State),
			URL:        fmt.Sprintf("https://bitbucket.org/%s/%s/issues/%d", owner, repo, issue.ID),
			Author:     issue.Reporter.Username,
			Labels:     []string{issue.Kind},
			CreatedAt:  issue.CreatedOn.Format(time.RFC3339),
			UpdatedAt:  issue.UpdatedOn.Format(time.RFC3339),
		}
	}
	return result, nil
}
