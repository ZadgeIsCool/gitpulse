package analysis

import (
	"fmt"
	"sort"
	"time"

	"github.com/ZadgeIsCool/gitpulse/github"
)

const staleIssueDays = 90

func analyzeIssues(client *github.Client, owner, repo string, since time.Time, report *Report) error {
	openIssues, err := client.ListIssues(owner, repo, "open")
	if err != nil {
		return fmt.Errorf("fetching open issues: %w", err)
	}

	closedIssues, err := client.ListIssues(owner, repo, "closed")
	if err != nil {
		return fmt.Errorf("fetching closed issues: %w", err)
	}

	report.Issues.OpenCount = len(openIssues)

	// Count closed in period
	closedInPeriod := 0
	for _, issue := range closedIssues {
		if issue.ClosedAt != nil && issue.ClosedAt.After(since) {
			closedInPeriod++
		}
	}
	report.Issues.ClosedCount = closedInPeriod

	// Close rate
	totalInPeriod := closedInPeriod + len(openIssues)
	if totalInPeriod > 0 {
		report.Issues.CloseRate = float64(closedInPeriod) / float64(totalInPeriod) * 100
	}

	// Stale issues: open issues with no update in staleIssueDays
	staleThreshold := time.Now().AddDate(0, 0, -staleIssueDays)
	report.Issues.StaleThresholdDays = staleIssueDays
	for _, issue := range openIssues {
		if issue.UpdatedAt.Before(staleThreshold) {
			report.Issues.StaleCount++
		}
	}

	// Median first response time (sample up to 30 recent issues)
	sampleSize := 30
	var recentIssues []github.Issue
	// Combine and sort by creation date
	allIssues := append(openIssues, closedIssues...)
	sort.Slice(allIssues, func(i, j int) bool {
		return allIssues[i].CreatedAt.After(allIssues[j].CreatedAt)
	})
	if len(allIssues) > sampleSize {
		recentIssues = allIssues[:sampleSize]
	} else {
		recentIssues = allIssues
	}

	var responseTimes []time.Duration
	for _, issue := range recentIssues {
		if issue.Comments == 0 {
			continue
		}
		comments, err := client.ListIssueComments(owner, repo, issue.Number)
		if err != nil {
			continue
		}
		for _, comment := range comments {
			// Skip if the commenter is the issue author
			if comment.User.Login != issue.User.Login {
				responseTimes = append(responseTimes, comment.CreatedAt.Sub(issue.CreatedAt))
				break
			}
		}
	}

	if len(responseTimes) > 0 {
		sort.Slice(responseTimes, func(i, j int) bool {
			return responseTimes[i] < responseTimes[j]
		})
		median := responseTimes[len(responseTimes)/2]
		report.Issues.MedianResponseHours = median.Hours()
	}

	return nil
}
