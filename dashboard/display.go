package dashboard

import (
	"fmt"
	"strings"

	"github.com/ZadgeIsCool/gitpulse/analysis"
)

// ANSI color codes
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	white     = "\033[37m"
	boldWhite = "\033[1;37m"
	bgBlue    = "\033[44m"
)

// Print renders the health dashboard to stdout.
func Print(r *analysis.Report) {
	width := 60

	printHeader(r, width)
	printHealthScore(r, width)
	printContributors(r, width)
	printPRMetrics(r, width)
	printIssues(r, width)
	printBusFactor(r, width)
	printDependencies(r, width)
	printFooter(r, width)
}

func printHeader(r *analysis.Report, w int) {
	line := strings.Repeat("─", w)
	fmt.Printf("\n%s%s╭%s╮%s\n", bold, cyan, line, reset)
	title := fmt.Sprintf("  GITPULSE — %s", r.Repo.FullName)
	padding := w - len(title) + 1
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s│%s%s%s%s│%s\n", bold, cyan, reset, boldWhite, title, strings.Repeat(" ", padding), reset)
	if r.Repo.Description != "" {
		desc := r.Repo.Description
		if len(desc) > w-4 {
			desc = desc[:w-7] + "..."
		}
		descPadding := w - len(desc) - 2
		if descPadding < 0 {
			descPadding = 0
		}
		fmt.Printf("%s%s│%s  %s%s%s%s│%s\n", bold, cyan, reset, dim, desc, strings.Repeat(" ", descPadding), cyan, reset)
	}
	info := fmt.Sprintf("  %s · %s★ %d%s · %s⑂ %d%s",
		colorIf(r.Repo.Language, blue),
		yellow, r.Repo.Stars, reset,
		green, r.Repo.Forks, reset)
	// We can't easily calculate visible length with ANSI, so use a simpler approach
	fmt.Printf("%s%s│%s%s\n", bold, cyan, reset, info)
	fmt.Printf("%s%s├%s┤%s\n", bold, cyan, strings.Repeat("─", w), reset)
}

func printHealthScore(r *analysis.Report, w int) {
	scoreColor := green
	scoreLabel := "Excellent"
	switch {
	case r.HealthScore >= 80:
		scoreColor = green
		scoreLabel = "Excellent"
	case r.HealthScore >= 60:
		scoreColor = green
		scoreLabel = "Good"
	case r.HealthScore >= 40:
		scoreColor = yellow
		scoreLabel = "Fair"
	case r.HealthScore >= 20:
		scoreColor = red
		scoreLabel = "Needs Attention"
	default:
		scoreColor = red
		scoreLabel = "Critical"
	}

	barLen := w - 20
	filled := r.HealthScore * barLen / 100
	empty := barLen - filled

	fmt.Printf("%s%s│%s %sHEALTH SCORE%s\n", bold, cyan, reset, boldWhite, reset)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
	fmt.Printf("%s%s│%s   %s%s%s%s %s%d/100%s — %s\n",
		bold, cyan, reset,
		scoreColor, strings.Repeat("█", filled), dim, strings.Repeat("░", empty),
		scoreColor, r.HealthScore, reset,
		scoreLabel)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
	printSeparator(w)
}

func printContributors(r *analysis.Report, w int) {
	fmt.Printf("%s%s│%s %sCONTRIBUTOR ACTIVITY%s (last %d days)\n",
		bold, cyan, reset, boldWhite, reset, r.AnalyzedDays)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)

	trendIcon := "→"
	trendColor := yellow
	switch r.Contributors.TrendDirection {
	case "up":
		trendIcon = "↑"
		trendColor = green
	case "down":
		trendIcon = "↓"
		trendColor = red
	}

	fmt.Printf("%s%s│%s   Total commits:       %s%d%s\n",
		bold, cyan, reset, boldWhite, r.Contributors.TotalCommits, reset)
	fmt.Printf("%s%s│%s   Active contributors: %s%d%s\n",
		bold, cyan, reset, boldWhite, r.Contributors.ActiveContribs, reset)
	fmt.Printf("%s%s│%s   Commits/week:        %s%.1f%s\n",
		bold, cyan, reset, boldWhite, r.Contributors.CommitsPerWeek, reset)
	fmt.Printf("%s%s│%s   Trend:               %s%s %.0f%%%s\n",
		bold, cyan, reset, trendColor, trendIcon, r.Contributors.TrendPercent, reset)

	if len(r.Contributors.TopContributors) > 0 {
		fmt.Printf("%s%s│%s\n", bold, cyan, reset)
		fmt.Printf("%s%s│%s   %sTop contributors:%s\n", bold, cyan, reset, dim, reset)
		limit := 5
		if len(r.Contributors.TopContributors) < limit {
			limit = len(r.Contributors.TopContributors)
		}
		for i := 0; i < limit; i++ {
			c := r.Contributors.TopContributors[i]
			barLen := int(c.Percent / 100 * 20)
			if barLen < 1 && c.Commits > 0 {
				barLen = 1
			}
			fmt.Printf("%s%s│%s   %s%-20s%s %s%s%s %d (%.0f%%)\n",
				bold, cyan, reset,
				white, c.Login, reset,
				magenta, strings.Repeat("▓", barLen), reset,
				c.Commits, c.Percent)
		}
	}

	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
	printSeparator(w)
}

func printPRMetrics(r *analysis.Report, w int) {
	fmt.Printf("%s%s│%s %sPULL REQUEST METRICS%s\n", bold, cyan, reset, boldWhite, reset)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)

	fmt.Printf("%s%s│%s   Open:     %s%d%s\n",
		bold, cyan, reset, boldWhite, r.PRMetrics.OpenCount, reset)
	fmt.Printf("%s%s│%s   Merged:   %s%d%s\n",
		bold, cyan, reset, green, r.PRMetrics.MergedCount, reset)
	fmt.Printf("%s%s│%s   Closed:   %s%d%s\n",
		bold, cyan, reset, red, r.PRMetrics.ClosedCount, reset)

	if r.PRMetrics.MedianReviewHours > 0 {
		reviewColor := green
		if r.PRMetrics.MedianReviewHours > 72 {
			reviewColor = red
		} else if r.PRMetrics.MedianReviewHours > 24 {
			reviewColor = yellow
		}
		fmt.Printf("%s%s│%s   Median review time:  %s%s%s\n",
			bold, cyan, reset, reviewColor, formatDuration(r.PRMetrics.MedianReviewHours), reset)
	}

	if r.PRMetrics.MedianMergeHours > 0 {
		mergeColor := green
		if r.PRMetrics.MedianMergeHours > 168 {
			mergeColor = red
		} else if r.PRMetrics.MedianMergeHours > 48 {
			mergeColor = yellow
		}
		fmt.Printf("%s%s│%s   Median merge time:   %s%s%s\n",
			bold, cyan, reset, mergeColor, formatDuration(r.PRMetrics.MedianMergeHours), reset)
	}

	if r.PRMetrics.StalePRs > 0 {
		fmt.Printf("%s%s│%s   Stale PRs:           %s%d%s (>30d old, no update in 14d)\n",
			bold, cyan, reset, red, r.PRMetrics.StalePRs, reset)
	}

	if r.PRMetrics.PRsWithoutReview > 0 {
		fmt.Printf("%s%s│%s   Merged without review: %s%d%s\n",
			bold, cyan, reset, yellow, r.PRMetrics.PRsWithoutReview, reset)
	}

	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
	printSeparator(w)
}

func printIssues(r *analysis.Report, w int) {
	fmt.Printf("%s%s│%s %sISSUE HEALTH%s\n", bold, cyan, reset, boldWhite, reset)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)

	fmt.Printf("%s%s│%s   Open issues:   %s%d%s\n",
		bold, cyan, reset, boldWhite, r.Issues.OpenCount, reset)
	fmt.Printf("%s%s│%s   Closed:        %s%d%s (in period)\n",
		bold, cyan, reset, green, r.Issues.ClosedCount, reset)

	if r.Issues.MedianResponseHours > 0 {
		respColor := green
		if r.Issues.MedianResponseHours > 168 {
			respColor = red
		} else if r.Issues.MedianResponseHours > 48 {
			respColor = yellow
		}
		fmt.Printf("%s%s│%s   Median response:  %s%s%s\n",
			bold, cyan, reset, respColor, formatDuration(r.Issues.MedianResponseHours), reset)
	}

	closeRateColor := green
	if r.Issues.CloseRate < 30 {
		closeRateColor = red
	} else if r.Issues.CloseRate < 50 {
		closeRateColor = yellow
	}
	fmt.Printf("%s%s│%s   Close rate:    %s%.0f%%%s\n",
		bold, cyan, reset, closeRateColor, r.Issues.CloseRate, reset)

	if r.Issues.StaleCount > 0 {
		fmt.Printf("%s%s│%s   Stale issues:  %s%d%s (no update in %dd)\n",
			bold, cyan, reset, red, r.Issues.StaleCount, reset, r.Issues.StaleThresholdDays)
	}

	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
	printSeparator(w)
}

func printBusFactor(r *analysis.Report, w int) {
	fmt.Printf("%s%s│%s %sBUS FACTOR%s\n", bold, cyan, reset, boldWhite, reset)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)

	bfColor := red
	switch {
	case r.BusFactor.BusFactor >= 5:
		bfColor = green
	case r.BusFactor.BusFactor >= 3:
		bfColor = green
	case r.BusFactor.BusFactor == 2:
		bfColor = yellow
	}

	fmt.Printf("%s%s│%s   Bus factor:    %s%d%s\n",
		bold, cyan, reset, bfColor, r.BusFactor.BusFactor, reset)
	fmt.Printf("%s%s│%s   Top contrib:   %.0f%% of commits\n",
		bold, cyan, reset, r.BusFactor.TopPercent)
	fmt.Printf("%s%s│%s   Assessment:    %s%s%s\n",
		bold, cyan, reset, bfColor, r.BusFactor.Assessment, reset)

	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
	printSeparator(w)
}

func printDependencies(r *analysis.Report, w int) {
	fmt.Printf("%s%s│%s %sDEPENDENCIES%s\n", bold, cyan, reset, boldWhite, reset)
	fmt.Printf("%s%s│%s\n", bold, cyan, reset)

	if !r.Dependencies.FileFound {
		fmt.Printf("%s%s│%s   %sNo recognized dependency file found%s\n",
			bold, cyan, reset, dim, reset)
	} else {
		fmt.Printf("%s%s│%s   Package manager: %s%s%s\n",
			bold, cyan, reset, boldWhite, r.Dependencies.FileType, reset)
		fmt.Printf("%s%s│%s   Dependencies:    %s%d%s\n",
			bold, cyan, reset, boldWhite, r.Dependencies.DepCount, reset)
	}

	fmt.Printf("%s%s│%s\n", bold, cyan, reset)
}

func printFooter(r *analysis.Report, w int) {
	line := strings.Repeat("─", w)
	fmt.Printf("%s%s╰%s╯%s\n", bold, cyan, line, reset)
	fmt.Printf("%s  Generated by gitpulse v0.1.0 · %s%s\n\n",
		dim, r.GeneratedAt.Format("2006-01-02 15:04:05"), reset)
}

func printSeparator(w int) {
	fmt.Printf("%s%s├%s┤%s\n", bold, cyan, strings.Repeat("─", w), reset)
}

func formatDuration(hours float64) string {
	if hours < 1 {
		return fmt.Sprintf("%.0f minutes", hours*60)
	}
	if hours < 48 {
		return fmt.Sprintf("%.1f hours", hours)
	}
	return fmt.Sprintf("%.1f days", hours/24)
}

func colorIf(s, color string) string {
	if s == "" {
		return dim + "unknown" + reset
	}
	return color + s + reset
}
