package analysis

// calculateHealthScore produces a 0-100 score based on all metrics.
func calculateHealthScore(r *Report) int {
	score := 0
	maxScore := 0

	// 1. Activity (25 points)
	maxScore += 25
	switch {
	case r.Contributors.CommitsPerWeek >= 10:
		score += 25
	case r.Contributors.CommitsPerWeek >= 5:
		score += 20
	case r.Contributors.CommitsPerWeek >= 2:
		score += 15
	case r.Contributors.CommitsPerWeek >= 0.5:
		score += 10
	case r.Contributors.TotalCommits > 0:
		score += 5
	}

	// Trend bonus/penalty (5 points)
	maxScore += 5
	switch r.Contributors.TrendDirection {
	case "up":
		score += 5
	case "stable":
		score += 3
	case "down":
		score += 1
	}

	// 2. PR health (20 points)
	maxScore += 20
	// Review time scoring
	if r.PRMetrics.MergedCount > 0 {
		switch {
		case r.PRMetrics.MedianReviewHours < 4:
			score += 10
		case r.PRMetrics.MedianReviewHours < 24:
			score += 7
		case r.PRMetrics.MedianReviewHours < 72:
			score += 4
		default:
			score += 1
		}
		// Low stale PR ratio
		if r.PRMetrics.OpenCount > 0 {
			staleRatio := float64(r.PRMetrics.StalePRs) / float64(r.PRMetrics.OpenCount)
			switch {
			case staleRatio < 0.1:
				score += 10
			case staleRatio < 0.3:
				score += 7
			case staleRatio < 0.5:
				score += 4
			default:
				score += 1
			}
		} else {
			score += 10
		}
	}

	// 3. Issue responsiveness (20 points)
	maxScore += 20
	if r.Issues.OpenCount+r.Issues.ClosedCount > 0 {
		// Response time
		switch {
		case r.Issues.MedianResponseHours > 0 && r.Issues.MedianResponseHours < 12:
			score += 10
		case r.Issues.MedianResponseHours < 48:
			score += 7
		case r.Issues.MedianResponseHours < 168:
			score += 4
		default:
			score += 1
		}
		// Close rate
		switch {
		case r.Issues.CloseRate > 70:
			score += 10
		case r.Issues.CloseRate > 50:
			score += 7
		case r.Issues.CloseRate > 30:
			score += 4
		default:
			score += 1
		}
	}

	// 4. Bus factor (20 points)
	maxScore += 20
	switch {
	case r.BusFactor.BusFactor >= 5:
		score += 20
	case r.BusFactor.BusFactor >= 3:
		score += 15
	case r.BusFactor.BusFactor == 2:
		score += 10
	case r.BusFactor.BusFactor == 1:
		score += 5
	}

	// 5. Not archived (10 points)
	maxScore += 10
	if !r.Repo.Archived {
		score += 10
	}

	if maxScore == 0 {
		return 0
	}
	return score * 100 / maxScore
}
