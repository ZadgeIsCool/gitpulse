package analysis

import (
	"sort"
	"time"

	"github.com/ZadgeIsCool/gitpulse/github"
)

func analyzeContributors(client *github.Client, owner, repo string, since time.Time, report *Report) error {
	commits, err := client.ListCommits(owner, repo, since)
	if err != nil {
		return err
	}

	report.Contributors.TotalCommits = len(commits)

	// Count commits per contributor
	contribMap := make(map[string]int)
	for _, c := range commits {
		login := ""
		if c.Author != nil {
			login = c.Author.Login
		}
		if login == "" {
			login = c.Commit.Author.Name
		}
		contribMap[login]++
	}

	report.Contributors.ActiveContribs = len(contribMap)

	// Sort by commit count
	type kv struct {
		Login   string
		Commits int
	}
	var sorted []kv
	for k, v := range contribMap {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Commits > sorted[j].Commits
	})

	// Top 10 contributors
	limit := 10
	if len(sorted) < limit {
		limit = len(sorted)
	}
	total := len(commits)
	for i := 0; i < limit; i++ {
		pct := 0.0
		if total > 0 {
			pct = float64(sorted[i].Commits) / float64(total) * 100
		}
		report.Contributors.TopContributors = append(report.Contributors.TopContributors, ContributorEntry{
			Login:   sorted[i].Login,
			Commits: sorted[i].Commits,
			Percent: pct,
		})
	}

	// Analyze trend: compare first half vs second half of period
	if len(commits) > 0 {
		midpoint := since.Add(time.Since(since) / 2)
		firstHalf, secondHalf := 0, 0
		for _, c := range commits {
			if c.Commit.Author.Date.Before(midpoint) {
				firstHalf++
			} else {
				secondHalf++
			}
		}

		weeks := float64(time.Since(since).Hours()) / (24 * 7)
		if weeks > 0 {
			report.Contributors.CommitsPerWeek = float64(total) / weeks
		}

		if firstHalf > 0 {
			change := (float64(secondHalf) - float64(firstHalf)) / float64(firstHalf) * 100
			report.Contributors.TrendPercent = change
			switch {
			case change > 10:
				report.Contributors.TrendDirection = "up"
			case change < -10:
				report.Contributors.TrendDirection = "down"
			default:
				report.Contributors.TrendDirection = "stable"
			}
		} else if secondHalf > 0 {
			report.Contributors.TrendDirection = "up"
			report.Contributors.TrendPercent = 100
		} else {
			report.Contributors.TrendDirection = "stable"
		}
	}

	return nil
}
