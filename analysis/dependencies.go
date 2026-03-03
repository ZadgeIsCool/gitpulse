package analysis

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/ZadgeIsCool/gitpulse/github"
)

type ghContent struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

var depFiles = []struct {
	path     string
	fileType string
}{
	{"go.mod", "Go modules"},
	{"package.json", "npm"},
	{"Cargo.toml", "Cargo (Rust)"},
	{"requirements.txt", "pip (Python)"},
	{"Gemfile", "Bundler (Ruby)"},
	{"pom.xml", "Maven (Java)"},
	{"build.gradle", "Gradle (Java)"},
	{"pyproject.toml", "Python"},
	{"Pipfile", "Pipenv (Python)"},
	{"composer.json", "Composer (PHP)"},
}

func analyzeDependencies(client *github.Client, owner, repo string, report *Report) {
	for _, df := range depFiles {
		data, err := client.GetDependencyFile(owner, repo, df.path)
		if err != nil || data == nil {
			continue
		}

		report.Dependencies.FileFound = true
		report.Dependencies.FileType = df.fileType

		// Decode content from GitHub API response
		var content ghContent
		if err := json.Unmarshal(data, &content); err != nil {
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(
			strings.ReplaceAll(content.Content, "\n", ""))
		if err != nil {
			continue
		}

		text := string(decoded)

		// Count dependencies based on file type
		switch df.path {
		case "go.mod":
			report.Dependencies.DepCount = countGoModDeps(text)
		case "package.json":
			report.Dependencies.DepCount = countPackageJSONDeps(text)
		case "requirements.txt":
			report.Dependencies.DepCount = countLineBasedDeps(text)
		case "Gemfile":
			report.Dependencies.DepCount = countGemfileDeps(text)
		default:
			report.Dependencies.DepCount = countLineBasedDeps(text)
		}

		return // Found a dependency file, stop looking
	}
}

func countGoModDeps(text string) int {
	count := 0
	inRequire := false
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		if inRequire && line != "" && !strings.HasPrefix(line, "//") {
			count++
		}
		if strings.HasPrefix(line, "require ") && !strings.Contains(line, "(") {
			count++
		}
	}
	return count
}

func countPackageJSONDeps(text string) int {
	var pkg map[string]json.RawMessage
	if err := json.Unmarshal([]byte(text), &pkg); err != nil {
		return 0
	}
	count := 0
	for _, key := range []string{"dependencies", "devDependencies"} {
		if raw, ok := pkg[key]; ok {
			var deps map[string]interface{}
			if json.Unmarshal(raw, &deps) == nil {
				count += len(deps)
			}
		}
	}
	return count
}

func countLineBasedDeps(text string) int {
	count := 0
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			count++
		}
	}
	return count
}

func countGemfileDeps(text string) int {
	count := 0
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "gem ") {
			count++
		}
	}
	return count
}
