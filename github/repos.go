package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetRepo fetches repository metadata.
func (c *Client) GetRepo(owner, repo string) (*Repo, error) {
	var r Repo
	err := c.get(fmt.Sprintf("/repos/%s/%s", owner, repo), &r)
	return &r, err
}

// ListCommits fetches commits since a given date.
func (c *Client) ListCommits(owner, repo string, since time.Time) ([]Commit, error) {
	var all []Commit
	page := 1
	for {
		var batch []Commit
		path := fmt.Sprintf("/repos/%s/%s/commits?since=%s&per_page=100&page=%d",
			owner, repo, since.Format(time.RFC3339), page)
		if err := c.get(path, &batch); err != nil {
			return all, err
		}
		all = append(all, batch...)
		if len(batch) < 100 {
			break
		}
		page++
		if page > 10 { // cap at 1000 commits
			break
		}
	}
	return all, nil
}

// ListPullRequests fetches recent pull requests (both open and closed).
func (c *Client) ListPullRequests(owner, repo string, state string) ([]PullRequest, error) {
	var all []PullRequest
	page := 1
	for {
		var batch []PullRequest
		path := fmt.Sprintf("/repos/%s/%s/pulls?state=%s&sort=updated&direction=desc&per_page=100&page=%d",
			owner, repo, state, page)
		if err := c.get(path, &batch); err != nil {
			return all, err
		}
		all = append(all, batch...)
		if len(batch) < 100 {
			break
		}
		page++
		if page > 5 {
			break
		}
	}
	return all, nil
}

// ListIssues fetches recent issues (excludes PRs via filtering).
func (c *Client) ListIssues(owner, repo string, state string) ([]Issue, error) {
	var all []Issue
	page := 1
	for {
		var batch []Issue
		path := fmt.Sprintf("/repos/%s/%s/issues?state=%s&sort=updated&direction=desc&per_page=100&page=%d",
			owner, repo, state, page)
		if err := c.get(path, &batch); err != nil {
			return all, err
		}
		// Filter out pull requests (GitHub returns PRs in the issues endpoint)
		for _, issue := range batch {
			if issue.PullRequestLinks == nil {
				all = append(all, issue)
			}
		}
		if len(batch) < 100 {
			break
		}
		page++
		if page > 5 {
			break
		}
	}
	return all, nil
}

// ListReviews fetches reviews for a specific pull request.
func (c *Client) ListReviews(owner, repo string, prNumber int) ([]Review, error) {
	var reviews []Review
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews?per_page=100", owner, repo, prNumber)
	err := c.get(path, &reviews)
	return reviews, err
}

// ListIssueComments fetches comments for a specific issue.
func (c *Client) ListIssueComments(owner, repo string, issueNumber int) ([]IssueComment, error) {
	var comments []IssueComment
	path := fmt.Sprintf("/repos/%s/%s/issues/%d/comments?per_page=1", owner, repo, issueNumber)
	err := c.get(path, &comments)
	return comments, err
}

// GetContributorStats fetches contributor statistics.
// This endpoint is cached by GitHub and may return 202 (computing).
func (c *Client) GetContributorStats(owner, repo string) ([]Contributor, error) {
	path := fmt.Sprintf("/repos/%s/%s/stats/contributors", owner, repo)

	for attempts := 0; attempts < 5; attempts++ {
		body, status, err := c.getRaw(path)
		if err != nil {
			return nil, err
		}
		if status == http.StatusAccepted {
			// GitHub is computing stats, wait and retry
			time.Sleep(2 * time.Second)
			continue
		}
		if status != http.StatusOK {
			return nil, fmt.Errorf("GitHub API error (status %d): %s", status, string(body))
		}
		var stats []Contributor
		if err := json.Unmarshal(body, &stats); err != nil {
			return nil, err
		}
		return stats, nil
	}
	return nil, fmt.Errorf("contributor stats not ready after retries (GitHub is still computing)")
}

// GetDependencyFile tries to fetch common dependency files to check freshness.
func (c *Client) GetDependencyFile(owner, repo, path string) ([]byte, error) {
	body, status, err := c.getRaw(fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path))
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s (status %d)", path, status)
	}
	return body, nil
}
