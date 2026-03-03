# gitpulse

A CLI tool that gives you a quick health check dashboard for any GitHub repository. Point it at a repo and get contributor activity trends, PR review times, issue response rates, stale issues, dependency info, and bus factor — all in one view.

Think of it as a fitness tracker for open source projects.

```
╭────────────────────────────────────────────────────────────╮
│  GITPULSE — golang/go                                      │
│  The Go programming language                               │
│  Go · ★ 125000 · ⑂ 17800                                  │
├────────────────────────────────────────────────────────────┤
│ HEALTH SCORE                                               │
│                                                            │
│   ██████████████████████████░░░░░░ 78/100 — Good           │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ CONTRIBUTOR ACTIVITY (last 90 days)                        │
│                                                            │
│   Total commits:       482                                 │
│   Active contributors: 87                                  │
│   Commits/week:        37.1                                │
│   Trend:               ↑ 12%                               │
│                                                            │
│   Top contributors:                                        │
│   gopherbot            ▓▓▓▓▓▓ 94 (20%)                    │
│   cherrymui            ▓▓▓ 52 (11%)                        │
│   ...                                                      │
├────────────────────────────────────────────────────────────┤
│ BUS FACTOR                                                 │
│                                                            │
│   Bus factor:    4                                         │
│   Top contrib:   20% of commits                            │
│   Assessment:    Reasonable — small but distributed team    │
╰────────────────────────────────────────────────────────────╯
```

## Install

### From source

Requires Go 1.21+.

```bash
git clone https://github.com/ZadgeIsCool/gitpulse.git
cd gitpulse
go build -o gitpulse .
```

Move the binary somewhere on your `$PATH`:

```bash
mv gitpulse /usr/local/bin/
```

### Cross-compile

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o gitpulse-linux .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o gitpulse-macos .

# Windows
GOOS=windows GOARCH=amd64 go build -o gitpulse.exe .
```

## Quick start

```bash
# Set your GitHub token
export GITHUB_TOKEN=ghp_your_token_here

# Analyze a repo
gitpulse golang/go

# Last 30 days only
gitpulse --days 30 rust-lang/rust

# JSON output for scripting
gitpulse --json kubernetes/kubernetes

# Pass token inline
gitpulse --token ghp_abc123 facebook/react
```

## Authentication

gitpulse requires a GitHub personal access token to access the API. You can provide it two ways:

| Method | Example |
|--------|---------|
| Environment variable | `export GITHUB_TOKEN=ghp_...` |
| CLI flag | `gitpulse --token ghp_... owner/repo` |

To create a token, go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens). For public repos, no special scopes are needed. For private repos, grant the `repo` scope.

## CLI reference

```
gitpulse [flags] owner/repo
```

| Flag | Default | Description |
|------|---------|-------------|
| `--token` | — | GitHub personal access token |
| `--days` | `90` | Number of days of history to analyze |
| `--json` | `false` | Output as JSON instead of the dashboard |
| `--version` | — | Print version and exit |

## What it measures

### Contributor activity

Analyzes all commits in the time window and produces:

- **Total commits** — raw count
- **Active contributors** — unique committers
- **Commits per week** — average rate
- **Trend** — compares the first half vs second half of the period. Marked as **↑ up** (>10% increase), **↓ down** (>10% decrease), or **→ stable**
- **Top contributors** — top 10 by commit count with percentage bars

### Pull request metrics

Fetches open and recently closed/merged PRs:

- **Open / Merged / Closed** — counts for the analysis period
- **Median review time** — time from PR creation to first substantive review (excludes author self-reviews and plain comments)
- **Median merge time** — time from PR creation to merge
- **Stale PRs** — open PRs older than 30 days with no update in the last 14 days
- **Merged without review** — PRs that were merged with no review from someone other than the author

Review time is sampled from up to 50 recently merged PRs to stay within API rate limits.

### Issue health

- **Open / Closed** — counts (closed filtered to the analysis period)
- **Median response time** — time from issue creation to the first comment by someone other than the author (sampled from 30 most recent issues)
- **Close rate** — percentage of issues closed vs total
- **Stale issues** — open issues with no update in 90+ days

### Bus factor

The minimum number of people whose combined commits account for 50% or more of all commits in the period.

| Bus factor | Assessment |
|-----------|------------|
| 1 (>80% from one person) | Critical — single point of failure |
| 1 | High risk — heavily dependent on one person |
| 2 | Moderate risk — small core team |
| 3–4 | Reasonable — small but distributed team |
| 5+ | Healthy — well-distributed contributions |

### Dependencies

Detects and parses the first recognized dependency file in the repo:

| File | Ecosystem |
|------|-----------|
| `go.mod` | Go modules |
| `package.json` | npm |
| `Cargo.toml` | Cargo (Rust) |
| `requirements.txt` | pip (Python) |
| `Gemfile` | Bundler (Ruby) |
| `pom.xml` | Maven (Java) |
| `build.gradle` | Gradle (Java) |
| `pyproject.toml` | Python |
| `Pipfile` | Pipenv (Python) |
| `composer.json` | Composer (PHP) |

Reports the package manager type and total dependency count.

## Health score

An overall score from 0 to 100, combining all metrics:

| Category | Points | What's measured |
|----------|--------|-----------------|
| **Activity** | 25 | Commits per week (0.5–10+ maps to 5–25 pts) |
| **Trend** | 5 | Up = 5, stable = 3, down = 1 |
| **PR health** | 20 | Review time (10 pts) + stale PR ratio (10 pts) |
| **Issue responsiveness** | 20 | Response time (10 pts) + close rate (10 pts) |
| **Bus factor** | 20 | 1 person = 5 pts, 5+ people = 20 pts |
| **Not archived** | 10 | 10 pts if the repo is active |

Score labels:

| Score | Label |
|-------|-------|
| 80–100 | Excellent |
| 60–79 | Good |
| 40–59 | Fair |
| 20–39 | Needs Attention |
| 0–19 | Critical |

## JSON output

Use `--json` for machine-readable output. The full schema:

```json
{
  "repo": {
    "full_name": "owner/repo",
    "description": "...",
    "stars": 1234,
    "forks": 567,
    "language": "Go",
    "archived": false
  },
  "contributors": {
    "total_commits": 482,
    "active_contributors": 87,
    "top_contributors": [
      { "login": "user1", "commits": 94, "percent": 19.5 }
    ],
    "commits_per_week": 37.1,
    "trend_direction": "up",
    "trend_percent": 12.3
  },
  "pull_requests": {
    "open_count": 45,
    "merged_count": 120,
    "closed_count": 15,
    "median_review_time_ns": 43200000000000,
    "median_review_hours": 12.0,
    "median_merge_time_ns": 86400000000000,
    "median_merge_hours": 24.0,
    "prs_without_review": 3,
    "stale_prs": 8
  },
  "issues": {
    "open_count": 200,
    "closed_count": 150,
    "median_response_hours": 6.5,
    "stale_issues": 42,
    "stale_threshold_days": 90,
    "close_rate_percent": 42.8
  },
  "bus_factor": {
    "bus_factor": 4,
    "top_contributor_percent": 19.5,
    "assessment": "Reasonable — small but distributed team"
  },
  "dependencies": {
    "file_found": true,
    "file_type": "Go modules",
    "dependency_count": 12,
    "last_modified": ""
  },
  "health_score": 78,
  "analyzed_days": 90,
  "generated_at": "2026-03-03T10:00:00Z"
}
```

## Project structure

```
gitpulse/
├── main.go                     Entry point
├── cmd/
│   └── root.go                 CLI flag parsing and dispatch
├── github/
│   ├── client.go               HTTP client with auth and pagination
│   ├── types.go                API response types
│   └── repos.go                Endpoint methods (commits, PRs, issues, reviews, deps)
├── analysis/
│   ├── report.go               Report struct and orchestration
│   ├── contributors.go         Commit activity and trend analysis
│   ├── pullrequests.go         PR review/merge time metrics
│   ├── issues.go               Issue response rates and staleness
│   ├── busfactor.go            Bus factor calculation
│   ├── dependencies.go         Dependency file detection and parsing
│   └── score.go                Health score calculation
├── dashboard/
│   ├── display.go              Colorized terminal output
│   └── json.go                 JSON output
├── go.mod
└── .gitignore
```

## How it works

1. Fetches repository metadata from the GitHub API
2. Pulls commits, PRs, and issues for the analysis window (paginated, capped for efficiency)
3. Samples merged PRs for review timing and recent issues for response timing
4. Calculates bus factor from contributor commit distribution
5. Checks for known dependency files and counts dependencies
6. Computes a weighted health score across all dimensions
7. Renders either a colorized terminal dashboard or JSON

API calls are made sequentially to respect GitHub rate limits. The contributor stats endpoint (`/stats/contributors`) is computed asynchronously by GitHub — gitpulse retries up to 5 times with 2-second delays if GitHub returns a 202 (computing) response.

## Rate limits

gitpulse makes roughly 80–120 API calls per analysis depending on the repo size. GitHub's rate limit is 5,000 requests/hour for authenticated users. You can comfortably analyze dozens of repos per hour.

## Zero dependencies

gitpulse uses only the Go standard library. No third-party packages. This means:

- Fast builds
- No supply chain risk
- Single static binary with no runtime dependencies
- Cross-compiles to any OS/architecture Go supports

## License

MIT
