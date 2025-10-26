package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// handleDashboard serves the HTML dashboard
func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "frontend/dashboard.html")
}

// handleGetCommits handles GET /commits?workspace={workspace}
func handleGetCommits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := r.URL.Query().Get("workspace")

	if workspace == "" {
		workspace = config.Workspace
	}

	// Get all repositories in the workspace
	repos, err := getRepositories(config.Email, config.Token, workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting repositories: %v", err), http.StatusInternalServerError)
		return
	}

	// Get commits for each repository
	allRepoCommits := []RepositoryCommits{}
	for _, repo := range repos {
		commits, err := getCommits(config.Email, config.Token, workspace, repo.Slug)
		if err != nil {
			log.Printf("Error getting commits for %s: %v", repo.Name, err)
			continue
		}

		repoCommits := RepositoryCommits{
			Repository: repo.FullName,
			Commits:    commits,
			Count:      len(commits),
		}
		allRepoCommits = append(allRepoCommits, repoCommits)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allRepoCommits)
}

// handleGetCommit handles GET /commit?workspace={workspace}&repo={repo}&hash={hash}
func handleGetCommit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := r.URL.Query().Get("workspace")
	repo := r.URL.Query().Get("repo")
	hash := r.URL.Query().Get("hash")

	if workspace == "" {
		workspace = config.Workspace
	}

	if repo == "" {
		http.Error(w, "repo parameter is required", http.StatusBadRequest)
		return
	}

	if hash == "" {
		http.Error(w, "hash parameter is required", http.StatusBadRequest)
		return
	}

	commit, err := getCommit(config.Email, config.Token, workspace, repo, hash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting commit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commit)
}

// handleGetPullRequests handles GET /pullrequests?workspace={workspace}&repo={repo}
func handleGetPullRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := r.URL.Query().Get("workspace")
	repo := r.URL.Query().Get("repo")

	if workspace == "" {
		workspace = config.Workspace
	}

	if repo == "" {
		http.Error(w, "repo parameter is required", http.StatusBadRequest)
		return
	}

	pullRequests, err := getPullRequests(config.Email, config.Token, workspace, repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting pull requests: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pullRequests)
}

// handleGetPullRequest handles GET /pullrequest?workspace={workspace}&repo={repo}&id={id}
func handleGetPullRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := r.URL.Query().Get("workspace")
	repo := r.URL.Query().Get("repo")
	idStr := r.URL.Query().Get("id")

	if workspace == "" {
		workspace = config.Workspace
	}

	if repo == "" {
		http.Error(w, "repo parameter is required", http.StatusBadRequest)
		return
	}

	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id must be a valid integer", http.StatusBadRequest)
		return
	}

	pullRequest, err := getPullRequest(config.Email, config.Token, workspace, repo, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting pull request: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pullRequest)
}

// handleGetRepositoryUsers handles GET /repository-users?workspace={workspace}
func handleGetRepositoryUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := r.URL.Query().Get("workspace")

	if workspace == "" {
		workspace = config.Workspace
	}

	// Get workspace members (all users in the workspace)
	workspaceUsers, err := getWorkspaceMembers(config.Email, config.Token, workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting workspace members: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaceUsers)
}

// handleGetUserCommits handles GET /user-commits?workspace={workspace}&account_id={account_id}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}
func handleGetUserCommits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Query().Get("account_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	workspace := config.Workspace

	if accountID == "" {
		http.Error(w, "account_id parameter is required", http.StatusBadRequest)
		return
	}

	// Get all repositories in the workspace
	repos, err := getRepositories(config.Email, config.Token, workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting repositories: %v", err), http.StatusInternalServerError)
		return
	}

	// Get commits for each repository and filter by account_id and date
	userCommits := []RepositoryCommits{}
	for _, repo := range repos {
		commits, err := getCommits(config.Email, config.Token, workspace, repo.Slug)
		if err != nil {
			log.Printf("Error getting commits for %s: %v", repo.Name, err)
			continue
		}

		// Filter commits by account_id and date range
		filteredCommits := []Commit{}
		for _, commit := range commits {
			if commit.Author.User.AccountID == accountID {
				// Apply date filter if provided
				if startDate != "" || endDate != "" {
					if isWithinDateRange(commit.Date, startDate, endDate) {
						filteredCommits = append(filteredCommits, commit)
					}
				} else {
					filteredCommits = append(filteredCommits, commit)
				}
			}
		}

		// Only include repositories where user has commits
		if len(filteredCommits) > 0 {
			repoCommits := RepositoryCommits{
				Repository: repo.FullName,
				Commits:    filteredCommits,
				Count:      len(filteredCommits),
			}
			userCommits = append(userCommits, repoCommits)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCommits)
}

// handleGetUserPullRequests handles GET /user-pullrequests?workspace={workspace}&account_id={account_id}
func handleGetUserPullRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Query().Get("account_id")
	workspace := config.Workspace

	if accountID == "" {
		http.Error(w, "account_id parameter is required", http.StatusBadRequest)
		return
	}

	// Get all repositories in the workspace
	repos, err := getRepositories(config.Email, config.Token, workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting repositories: %v", err), http.StatusInternalServerError)
		return
	}

	// Get pull requests for each repository and filter by account_id
	userPRs := []RepositoryPullRequests{}
	for _, repo := range repos {
		prs, err := getPullRequests(config.Email, config.Token, workspace, repo.Slug)
		if err != nil {
			log.Printf("Error getting pull requests for %s: %v", repo.Name, err)
			continue
		}

		// Filter PRs by author account_id (merged PRs only)
		filteredPRs := []PullRequest{}
		for _, pr := range prs {
			if pr.Author.AccountID == accountID {
				filteredPRs = append(filteredPRs, pr)
			}
		}

		// Only include repositories where user has merged PRs
		if len(filteredPRs) > 0 {
			repoPRs := RepositoryPullRequests{
				Repository:   repo.FullName,
				PullRequests: filteredPRs,
				Count:        len(filteredPRs),
			}
			userPRs = append(userPRs, repoPRs)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userPRs)
}

// handleGetUserCommitFrequency handles GET /user-commit-frequency?workspace={workspace}&account_id={account_id}
func handleGetUserCommitFrequency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Query().Get("account_id")
	workspace := config.Workspace

	if accountID == "" {
		http.Error(w, "account_id parameter is required", http.StatusBadRequest)
		return
	}

	// Get all repositories in the workspace
	repos, err := getRepositories(config.Email, config.Token, workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting repositories: %v", err), http.StatusInternalServerError)
		return
	}

	// Collect all commits by the user across all repositories
	var allCommits []Commit
	commitsByRepo := make(map[string]int)

	for _, repo := range repos {
		commits, err := getCommits(config.Email, config.Token, workspace, repo.Slug)
		if err != nil {
			log.Printf("Error getting commits for %s: %v", repo.Name, err)
			continue
		}

		// Filter commits by account_id
		for _, commit := range commits {
			if commit.Author.User.AccountID == accountID {
				allCommits = append(allCommits, commit)
				commitsByRepo[repo.FullName]++
			}
		}
	}

	// Calculate frequency statistics
	frequency := calculateCommitFrequency(allCommits, commitsByRepo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(frequency)
}

// handleGetLanguages handles GET /languages?workspace={workspace}&repo={repo}
func handleGetLanguages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspace := r.URL.Query().Get("workspace")
	repo := r.URL.Query().Get("repo")

	if workspace == "" {
		workspace = config.Workspace
	}

	// If specific repo is requested
	if repo != "" {
		languages, err := getRepositoryLanguages(config.Email, config.Token, workspace, repo)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting languages for repository: %v", err), http.StatusInternalServerError)
			return
		}

		repoLanguages := RepositoryLanguages{
			Repository: fmt.Sprintf("%s/%s", workspace, repo),
			Languages:  languages,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repoLanguages)
		return
	}

	// Get languages for all repositories in workspace
	repos, err := getRepositories(config.Email, config.Token, workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting repositories: %v", err), http.StatusInternalServerError)
		return
	}

	allRepoLanguages := []RepositoryLanguages{}
	for _, repo := range repos {
		languages, err := getRepositoryLanguages(config.Email, config.Token, workspace, repo.Slug)
		if err != nil {
			log.Printf("Error getting languages for %s: %v", repo.Name, err)
			continue
		}

		// Only include if languages were found
		if len(languages) > 0 {
			repoLanguages := RepositoryLanguages{
				Repository: repo.FullName,
				Languages:  languages,
			}
			allRepoLanguages = append(allRepoLanguages, repoLanguages)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allRepoLanguages)
}
