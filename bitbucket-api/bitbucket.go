package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const bitbucketAPIBase = "https://api.bitbucket.org/2.0"

// getRepositories fetches all repositories for a workspace
func getRepositories(email, token, workspace string) ([]Repository, error) {
	url := fmt.Sprintf("%s/repositories/%s", bitbucketAPIBase, workspace)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var reposResp RepositoriesResponse
	if err := json.Unmarshal(body, &reposResp); err != nil {
		return nil, fmt.Errorf("error parsing repositories JSON: %v", err)
	}

	return reposResp.Values, nil
}

// getCommits fetches all commits for a repository
func getCommits(email, token, workspace, repoSlug string) ([]Commit, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/commits", bitbucketAPIBase, workspace, repoSlug)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var commitsResp CommitsResponse
	if err := json.Unmarshal(body, &commitsResp); err != nil {
		return nil, fmt.Errorf("error parsing commits JSON: %v", err)
	}

	return commitsResp.Values, nil
}

// getCommit fetches a specific commit by hash
func getCommit(email, token, workspace, repoSlug, hash string) (*Commit, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/commit/%s", bitbucketAPIBase, workspace, repoSlug, hash)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var commit Commit
	if err := json.Unmarshal(body, &commit); err != nil {
		return nil, fmt.Errorf("error parsing commit JSON: %v", err)
	}

	return &commit, nil
}

// getPullRequests fetches all merged pull requests for a repository
func getPullRequests(email, token, workspace, repoSlug string) ([]PullRequest, error) {
	// Get merged pull requests with fields to include email
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?fields=values.id,values.title,values.state,values.author.display_name,values.author.nickname,values.author.account_id,values.author.uuid,values.closed_by.display_name,values.closed_by.nickname,values.closed_by.account_id,values.closed_by.uuid,values.created_on,values.updated_on,size", bitbucketAPIBase, workspace, repoSlug)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var prResp PullRequestsResponse
	if err := json.Unmarshal(body, &prResp); err != nil {
		return nil, fmt.Errorf("error parsing pull requests JSON: %v", err)
	}

	return prResp.Values, nil
}

// getPullRequest fetches a specific pull request by ID
func getPullRequest(email, token, workspace, repoSlug string, prID int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d", bitbucketAPIBase, workspace, repoSlug, prID)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var pr PullRequest
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, fmt.Errorf("error parsing pull request JSON: %v", err)
	}

	return &pr, nil
}

// getWorkspaceMembers fetches all members of a workspace
func getWorkspaceMembers(email, token, workspace string) ([]User, error) {
	url := fmt.Sprintf("%s/workspaces/%s/members", bitbucketAPIBase, workspace)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var membersResp WorkspaceMembersResponse
	if err := json.Unmarshal(body, &membersResp); err != nil {
		return nil, fmt.Errorf("error parsing workspace members JSON: %v", err)
	}

	users := make([]User, 0, len(membersResp.Values))
	for _, member := range membersResp.Values {
		users = append(users, member.User)
	}

	return users, nil
}

// makeRequest performs an authenticated HTTP GET request to the Bitbucket API
func makeRequest(email, token, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.SetBasicAuth(email, token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// isWithinDateRange checks if a commit date is within the specified range
// Dates should be in format YYYY-MM-DD
// commitDate is in ISO 8601 format from Bitbucket API
func isWithinDateRange(commitDate, startDate, endDate string) bool {
	if commitDate == "" {
		return false
	}

	// Extract date part from ISO 8601 timestamp (YYYY-MM-DDTHH:MM:SS+00:00 -> YYYY-MM-DD)
	if len(commitDate) >= 10 {
		commitDate = commitDate[:10]
	}

	// If start_date is provided, check if commit is on or after start_date
	if startDate != "" && commitDate < startDate {
		return false
	}

	// If end_date is provided, check if commit is on or before end_date
	if endDate != "" && commitDate > endDate {
		return false
	}

	return true
}

// calculateCommitFrequency calculates commit frequency statistics
func calculateCommitFrequency(commits []Commit, commitsByRepo map[string]int) CommitFrequency {
	if len(commits) == 0 {
		return CommitFrequency{
			TotalCommits:       0,
			DateRange:          "No commits found",
			AveragePerDay:      0,
			AveragePerWeek:     0,
			MostActiveDay:      "N/A",
			CommitsByDay:       make(map[string]int),
			CommitsByDayOfWeek: make(map[string]int),
			CommitsByRepo:      commitsByRepo,
		}
	}

	commitsByDay := make(map[string]int)
	commitsByDayOfWeek := make(map[string]int)
	var firstDate, lastDate string

	dayOfWeekNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	for _, commit := range commits {
		if commit.Date == "" {
			continue
		}

		// Extract date (YYYY-MM-DD)
		date := commit.Date
		if len(date) >= 10 {
			date = date[:10]
		}

		// Track first and last dates
		if firstDate == "" || date < firstDate {
			firstDate = date
		}
		if lastDate == "" || date > lastDate {
			lastDate = date
		}

		// Count by day
		commitsByDay[date]++

		// Count by day of week
		// Parse the date to get day of week
		if t, err := parseDate(date); err == nil {
			dayName := dayOfWeekNames[t.Weekday()]
			commitsByDayOfWeek[dayName]++
		}
	}

	// Calculate date range
	dateRange := fmt.Sprintf("%s to %s", firstDate, lastDate)

	// Calculate number of days
	days := calculateDaysBetween(firstDate, lastDate)
	if days == 0 {
		days = 1
	}

	// Calculate averages
	averagePerDay := float64(len(commits)) / float64(days)
	averagePerWeek := averagePerDay * 7

	// Find most active day
	mostActiveDay := ""
	maxCommits := 0
	for day, count := range commitsByDay {
		if count > maxCommits {
			maxCommits = count
			mostActiveDay = day
		}
	}

	return CommitFrequency{
		TotalCommits:       len(commits),
		DateRange:          dateRange,
		AveragePerDay:      averagePerDay,
		AveragePerWeek:     averagePerWeek,
		MostActiveDay:      mostActiveDay,
		CommitsByDay:       commitsByDay,
		CommitsByDayOfWeek: commitsByDayOfWeek,
		CommitsByRepo:      commitsByRepo,
	}
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// calculateDaysBetween calculates the number of days between two dates
func calculateDaysBetween(startDate, endDate string) int {
	start, err1 := parseDate(startDate)
	end, err2 := parseDate(endDate)

	if err1 != nil || err2 != nil {
		return 1
	}

	diff := end.Sub(start)
	days := int(diff.Hours() / 24)
	if days < 1 {
		days = 1
	}

	return days
}

// extensionToLanguage maps file extensions to programming languages
var extensionToLanguage = map[string]string{
	".go":         "Go",
	".js":         "JavaScript",
	".ts":         "TypeScript",
	".tsx":        "TypeScript",
	".jsx":        "JavaScript",
	".py":         "Python",
	".java":       "Java",
	".c":          "C",
	".cpp":        "C++",
	".cc":         "C++",
	".cxx":        "C++",
	".h":          "C/C++ Header",
	".hpp":        "C++ Header",
	".cs":         "C#",
	".php":        "PHP",
	".rb":         "Ruby",
	".rs":         "Rust",
	".swift":      "Swift",
	".kt":         "Kotlin",
	".scala":      "Scala",
	".html":       "HTML",
	".css":        "CSS",
	".scss":       "SCSS",
	".sass":       "Sass",
	".less":       "Less",
	".vue":        "Vue",
	".sql":        "SQL",
	".sh":         "Shell",
	".bash":       "Shell",
	".yaml":       "YAML",
	".yml":        "YAML",
	".json":       "JSON",
	".xml":        "XML",
	".r":          "R",
	".m":          "Objective-C",
	".dart":       "Dart",
	".lua":        "Lua",
	".pl":         "Perl",
	".groovy":     "Groovy",
	".dockerfile": "Dockerfile",
}

// getRepositoryFileTree fetches the file tree for a repository
func getRepositoryFileTree(email, token, workspace, repoSlug string) ([]FileNode, error) {
	// Use HEAD to get the default branch's file tree
	url := fmt.Sprintf("%s/repositories/%s/%s/src/HEAD/", bitbucketAPIBase, workspace, repoSlug)

	body, err := makeRequest(email, token, url)
	if err != nil {
		return nil, err
	}

	var fileTreeResp FileTreeResponse
	if err := json.Unmarshal(body, &fileTreeResp); err != nil {
		return nil, fmt.Errorf("error parsing file tree JSON: %v", err)
	}

	return fileTreeResp.Values, nil
}

// analyzeLanguages analyzes file extensions and calculates language percentages
func analyzeLanguages(files []FileNode) []Language {
	languageCount := make(map[string]int)
	totalFiles := 0

	// Count files by language
	for _, file := range files {
		if file.Type != "commit_file" {
			continue
		}

		// Get file extension
		ext := ""
		for i := len(file.Path) - 1; i >= 0; i-- {
			if file.Path[i] == '.' {
				ext = file.Path[i:]
				break
			}
			if file.Path[i] == '/' {
				break
			}
		}

		// Map extension to language
		if language, ok := extensionToLanguage[ext]; ok {
			languageCount[language]++
			totalFiles++
		}
	}

	// Convert to Language structs and calculate percentages
	languages := []Language{}
	for name, count := range languageCount {
		percentage := 0.0
		if totalFiles > 0 {
			percentage = (float64(count) / float64(totalFiles)) * 100
		}
		languages = append(languages, Language{
			Name:       name,
			Percentage: percentage,
			FileCount:  count,
		})
	}

	// Sort by percentage (descending)
	for i := 0; i < len(languages); i++ {
		for j := i + 1; j < len(languages); j++ {
			if languages[j].Percentage > languages[i].Percentage {
				languages[i], languages[j] = languages[j], languages[i]
			}
		}
	}

	return languages
}

// getRepositoryLanguages fetches and analyzes language statistics for a repository
func getRepositoryLanguages(email, token, workspace, repoSlug string) ([]Language, error) {
	files, err := getRepositoryFileTree(email, token, workspace, repoSlug)
	if err != nil {
		return nil, err
	}

	return analyzeLanguages(files), nil
}
