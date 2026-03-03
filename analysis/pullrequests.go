package analysis

import (
	"fmt"
	"sort"
	"time"

	"github.com/ZadgeIsCool/gitpulse/github"
)

func analyzePRs(client *github.Client, owner, repo string, since time.Time, report *Report) error {
	// Fetch both open and closed PRs
	openPRs, err := client.ListPullRequests(owner, repo, "open")
	if err != nil {
		return fmt.Errorf("fetching open PRs: %w", err)
	}

	closedPRs, err := client.ListPullRequests(owner, repo, "closed")
	if err != nil {
		return fmt.Errorf("fetching closed PRs: %w", err)
	}

	report.PRMetrics.OpenCount = len(openPRs)

	// Filter closed PRs to analysis window
	var mergedInPeriod []github.PullRequest
	var closedInPeriod []github.PullRequest
	for _, pr := range closedPRs {
		if pr.ClosedAt != nil && pr.ClosedAt.After(since) {
			if pr.MergedAt != nil {
				mergedInPeriod = append(mergedInPeriod, pr)
			} else {
				closedInPeriod = append(closedInPeriod, pr)
			}
		}
	}

	report.PRMetrics.MergedCount = len(mergedInPeriod)
	report.PRMetrics.ClosedCount = len(closedInPeriod)

	// Calculate median time to first review for merged PRs
	var reviewTimes []time.Duration
	// Sample up to 50 merged PRs for review times (to avoid too many API calls)
	sampleSize := 50
	if len(mergedInPeriod) < sampleSize {
		sampleSize = len(mergedInPeriod)
	}
	noReviewCount := 0
	for i := 0; i < sampleSize; i++ {
		pr := mergedInPeriod[i]
		reviews, err := client.ListReviews(owner, repo, pr.Number)
		if err != nil {
			continue
		}
		if len(reviews) == 0 {
			noReviewCount++
			continue
		}
		// Find first non-author review
		var firstReview *time.Time
		for _, r := range reviews {
			if r.User.Login != pr.User.Login && r.State != "COMMENTED" {
				t := r.SubmittedAt
				if firstReview == nil || t.Before(*firstReview) {
					firstReview = &t
				}
			}
		}
		if firstReview != nil {
			reviewTimes = append(reviewTimes, firstReview.Sub(pr.CreatedAt))
		} else {
			noReviewCount++
		}
	}

	report.PRMetrics.PRsWithoutReview = noReviewCount

	if len(reviewTimes) > 0 {
		sort.Slice(reviewTimes, func(i, j int) bool {
			return reviewTimes[i] < reviewTimes[j]
		})
		median := reviewTimes[len(reviewTimes)/2]
		report.PRMetrics.MedianReviewTime = median
		report.PRMetrics.MedianReviewHours = median.Hours()
	}

	// Median merge time
	var mergeTimes []time.Duration
	for _, pr := range mergedInPeriod {
		if pr.MergedAt != nil {
			mergeTimes = append(mergeTimes, pr.MergedAt.Sub(pr.CreatedAt))
		}
	}
	if len(mergeTimes) > 0 {
		sort.Slice(mergeTimes, func(i, j int) bool {
			return mergeTimes[i] < mergeTimes[j]
		})
		median := mergeTimes[len(mergeTimes)/2]
		report.PRMetrics.MedianMergeTime = median
		report.PRMetrics.MedianMergeHours = median.Hours()
	}

	// Stale PRs: open for more than 30 days with no update in 14 days
	staleThreshold := time.Now().AddDate(0, 0, -14)
	for _, pr := range openPRs {
		if pr.UpdatedAt.Before(staleThreshold) && time.Since(pr.CreatedAt).Hours() > 30*24 {
			report.PRMetrics.StalePRs++
		}
	}

	return nil
}
