package analysis

func analyzeBusFactor(report *Report) {
	contribs := report.Contributors.TopContributors
	if len(contribs) == 0 {
		report.BusFactor = BusFactorAnalysis{
			BusFactor:  0,
			Assessment: "No data",
		}
		return
	}

	totalCommits := report.Contributors.TotalCommits
	if totalCommits == 0 {
		report.BusFactor = BusFactorAnalysis{
			BusFactor:  0,
			Assessment: "No commits in period",
		}
		return
	}

	// Bus factor = minimum number of people whose contributions
	// account for >= 50% of all commits
	cumulative := 0
	busFactor := 0
	for _, c := range contribs {
		cumulative += c.Commits
		busFactor++
		if float64(cumulative)/float64(totalCommits) >= 0.5 {
			break
		}
	}

	topPercent := contribs[0].Percent

	var assessment string
	switch {
	case busFactor == 1 && topPercent > 80:
		assessment = "Critical — single point of failure"
	case busFactor == 1:
		assessment = "High risk — heavily dependent on one person"
	case busFactor == 2:
		assessment = "Moderate risk — small core team"
	case busFactor <= 4:
		assessment = "Reasonable — small but distributed team"
	default:
		assessment = "Healthy — well-distributed contributions"
	}

	report.BusFactor = BusFactorAnalysis{
		BusFactor:  busFactor,
		TopPercent: topPercent,
		Assessment: assessment,
	}
}
