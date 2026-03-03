package analysis

import (
	"time"

	"github.com/ZadgeIsCool/gitpulse/github"
)

// Report is the full health check result.
type Report struct {
	Repo         RepoSummary         `json:"repo"`
	Contributors ContributorAnalysis `json:"contributors"`
	PRMetrics    PRAnalysis          `json:"pull_requests"`
	Issues       IssueAnalysis       `json:"issues"`
	BusFactor    BusFactorAnalysis   `json:"bus_factor"`
	Dependencies DependencyAnalysis  `json:"dependencies"`
	HealthScore  int                 `json:"health_score"`
	AnalyzedDays int                 `json:"analyzed_days"`
	GeneratedAt  time.Time           `json:"generated_at"`
}

type RepoSummary struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Language    string `json:"language"`
	Archived    bool   `json:"archived"`
}

type ContributorAnalysis struct {
	TotalCommits     int                `json:"total_commits"`
	ActiveContribs   int                `json:"active_contributors"`
	TopContributors  []ContributorEntry `json:"top_contributors"`
	CommitsPerWeek   float64            `json:"commits_per_week"`
	TrendDirection   string             `json:"trend_direction"` // "up", "down", "stable"
	TrendPercent     float64            `json:"trend_percent"`
}

type ContributorEntry struct {
	Login   string `json:"login"`
	Commits int    `json:"commits"`
	Percent float64 `json:"percent"`
}

type PRAnalysis struct {
	OpenCount          int           `json:"open_count"`
	MergedCount        int           `json:"merged_count"`
	ClosedCount        int           `json:"closed_count"`
	MedianReviewTime   time.Duration `json:"median_review_time_ns"`
	MedianReviewHours  float64       `json:"median_review_hours"`
	MedianMergeTime    time.Duration `json:"median_merge_time_ns"`
	MedianMergeHours   float64       `json:"median_merge_hours"`
	PRsWithoutReview   int           `json:"prs_without_review"`
	StalePRs           int           `json:"stale_prs"`
}

type IssueAnalysis struct {
	OpenCount           int     `json:"open_count"`
	ClosedCount         int     `json:"closed_count"`
	MedianResponseHours float64 `json:"median_response_hours"`
	StaleCount          int     `json:"stale_issues"`
	StaleThresholdDays  int     `json:"stale_threshold_days"`
	CloseRate           float64 `json:"close_rate_percent"`
}

type BusFactorAnalysis struct {
	BusFactor     int     `json:"bus_factor"`
	TopPercent    float64 `json:"top_contributor_percent"`
	Assessment    string  `json:"assessment"`
}

type DependencyAnalysis struct {
	FileFound    bool   `json:"file_found"`
	FileType     string `json:"file_type"`
	DepCount     int    `json:"dependency_count"`
	LastModified string `json:"last_modified"`
}

// GenerateReport runs all analysis modules and produces a full report.
func GenerateReport(client *github.Client, owner, repo string, days int) (*Report, error) {
	since := time.Now().AddDate(0, 0, -days)

	// Fetch repo info
	repoInfo, err := client.GetRepo(owner, repo)
	if err != nil {
		return nil, err
	}

	report := &Report{
		Repo: RepoSummary{
			FullName:    repoInfo.FullName,
			Description: repoInfo.Description,
			Stars:       repoInfo.Stars,
			Forks:       repoInfo.Forks,
			Language:    repoInfo.Language,
			Archived:    repoInfo.Archived,
		},
		AnalyzedDays: days,
		GeneratedAt:  time.Now(),
	}

	// Run analysis modules (sequential to respect rate limits)
	if err := analyzeContributors(client, owner, repo, since, report); err != nil {
		return nil, err
	}

	if err := analyzePRs(client, owner, repo, since, report); err != nil {
		return nil, err
	}

	if err := analyzeIssues(client, owner, repo, since, report); err != nil {
		return nil, err
	}

	analyzeBusFactor(report)
	analyzeDependencies(client, owner, repo, report)

	report.HealthScore = calculateHealthScore(report)

	return report, nil
}
