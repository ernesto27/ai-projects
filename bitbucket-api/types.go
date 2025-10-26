package main

// Repository represents a Bitbucket repository
type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Slug     string `json:"slug"`
}

// RepositoriesResponse represents the API response for repositories
type RepositoriesResponse struct {
	Values []Repository `json:"values"`
}

// Author represents a commit author
type Author struct {
	Raw  string `json:"raw"`
	User struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
		AccountID   string `json:"account_id"`
	} `json:"user"`
}

// Commit represents a single commit
type Commit struct {
	Hash   string `json:"hash"`
	Author Author `json:"author"`
	Date   string `json:"date"`
}

// CommitsResponse represents the API response for commits
type CommitsResponse struct {
	Values []Commit `json:"values"`
	Size   int      `json:"size"`
}

// RepositoryCommits represents commits grouped by repository
type RepositoryCommits struct {
	Repository string   `json:"repository"`
	Commits    []Commit `json:"commits"`
	Count      int      `json:"count"`
}

// User represents a Bitbucket user
type User struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Nickname    string `json:"nickname"`
	AccountID   string `json:"account_id"`
}

// PullRequest represents a Bitbucket pull request
type PullRequest struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Author struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
		AccountID   string `json:"account_id"`
	} `json:"author"`
	ClosedBy *struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
		AccountID   string `json:"account_id"`
	} `json:"closed_by"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

// PullRequestsResponse represents the API response for pull requests
type PullRequestsResponse struct {
	Values []PullRequest `json:"values"`
	Size   int           `json:"size"`
}

// RepositoryUsers represents users grouped by repository
type RepositoryUsers struct {
	Repository string `json:"repository"`
	Users      []User `json:"users"`
	Count      int    `json:"count"`
}

// WorkspaceMember represents a workspace member
type WorkspaceMember struct {
	User User `json:"user"`
}

// WorkspaceMembersResponse represents the API response for workspace members
type WorkspaceMembersResponse struct {
	Values []WorkspaceMember `json:"values"`
}

// UserCommitsSummary represents a summary of user commits
type UserCommitsSummary struct {
	User              User                `json:"user"`
	TotalCommits      int                 `json:"total_commits"`
	RepositoryCommits []RepositoryCommits `json:"repository_commits"`
}

// RepositoryPullRequests represents pull requests grouped by repository
type RepositoryPullRequests struct {
	Repository   string        `json:"repository"`
	PullRequests []PullRequest `json:"pull_requests"`
	Count        int           `json:"count"`
}

// CommitFrequency represents commit frequency statistics
type CommitFrequency struct {
	TotalCommits       int            `json:"total_commits"`
	DateRange          string         `json:"date_range"`
	AveragePerDay      float64        `json:"average_per_day"`
	AveragePerWeek     float64        `json:"average_per_week"`
	MostActiveDay      string         `json:"most_active_day"`
	CommitsByDay       map[string]int `json:"commits_by_day"`
	CommitsByDayOfWeek map[string]int `json:"commits_by_day_of_week"`
	CommitsByRepo      map[string]int `json:"commits_by_repo"`
}

// Config holds the Bitbucket API credentials
type Config struct {
	Email     string
	Token     string
	Workspace string
}

// Language represents programming language usage in a repository
type Language struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	FileCount  int     `json:"file_count"`
}

// RepositoryLanguages represents languages used in a repository
type RepositoryLanguages struct {
	Repository string     `json:"repository"`
	Languages  []Language `json:"languages"`
}

// FileTreeResponse represents the API response for repository file tree
type FileTreeResponse struct {
	Values []FileNode `json:"values"`
}

// FileNode represents a file or directory in the repository
type FileNode struct {
	Path string `json:"path"`
	Type string `json:"type"` // "commit_file" or "commit_directory"
}
