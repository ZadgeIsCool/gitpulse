package github

import "time"

// Repo represents a GitHub repository.
type Repo struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	Stars           int       `json:"stargazers_count"`
	Forks           int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Language        string    `json:"language"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Archived        bool      `json:"archived"`
	DefaultBranch   string    `json:"default_branch"`
}

// Commit represents a GitHub commit.
type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	Author *User `json:"author"`
}

// User represents a GitHub user.
type User struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at"`
	MergedAt  *time.Time `json:"merged_at"`
	User      User      `json:"user"`
	Draft     bool      `json:"draft"`
}

// Issue represents a GitHub issue (not a PR).
type Issue struct {
	Number           int        `json:"number"`
	Title            string     `json:"title"`
	State            string     `json:"state"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	User             User       `json:"user"`
	Comments         int        `json:"comments"`
	PullRequestLinks *struct{}  `json:"pull_request"`
}

// IssueComment represents a comment on an issue.
type IssueComment struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
}

// Review represents a pull request review.
type Review struct {
	ID          int       `json:"id"`
	User        User      `json:"user"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// Contributor represents contributor stats.
type Contributor struct {
	Author User `json:"author"`
	Total  int  `json:"total"`
	Weeks  []struct {
		W int `json:"w"` // unix timestamp of week start
		A int `json:"a"` // additions
		D int `json:"d"` // deletions
		C int `json:"c"` // commits
	} `json:"weeks"`
}
