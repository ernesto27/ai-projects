# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Bitbucket Analytics API and Dashboard - a Go web service that fetches repository data from the Bitbucket API and presents analytics through both a REST API and an interactive HTML dashboard.

## Architecture

### Application Structure

The codebase follows a simple, flat architecture with four main Go files:

- **main.go**: Entry point, HTTP server setup, and route registration. Loads configuration from environment variables (BITBUCKET_EMAIL, BITBUCKET_TOKEN, BITBUCKET_WORKSPACE) and starts the HTTP server on port 8080 (configurable via PORT env var).

- **bitbucket.go**: Bitbucket API client layer. Contains all functions that make authenticated requests to the Bitbucket API v2.0. Key functions include:
  - `getRepositories()`: Fetches all repos in a workspace
  - `getCommits()`: Fetches commits for a specific repo
  - `getCommit()`: Fetches a specific commit by hash
  - `getPullRequests()`: Fetches merged PRs with author details
  - `getWorkspaceMembers()`: Fetches workspace users
  - `makeRequest()`: Central HTTP client with Basic Auth
  - `calculateCommitFrequency()`: Computes commit statistics (averages, date ranges, day-of-week distribution)

- **controllers.go**: HTTP handlers for all endpoints. Each handler validates query params, calls appropriate bitbucket.go functions, aggregates data across repos, and returns JSON responses.

- **types.go**: Data structures for API requests/responses. Includes types for Repository, Commit, PullRequest, User, and various response wrappers. Note the nested JSON structure for authors (commit.Author.User.AccountID).

### Data Flow

1. User requests data via dashboard or API endpoint
2. Controller validates params and extracts workspace/repo/user identifiers
3. Controller calls bitbucket.go functions to fetch raw data from Bitbucket API
4. For aggregate endpoints (e.g., /user-commits), controller iterates over all repos in workspace
5. Data is filtered/transformed (e.g., by account_id, date range) and returned as JSON

### Key Patterns

- **User Identification**: Users are tracked by `account_id` (unique identifier). Display names and nicknames are used for UI but account_id is the stable key for filtering.
- **Workspace-scoped Operations**: Most endpoints iterate over ALL repositories in a workspace to aggregate data (commits, PRs, etc.)
- **Date Filtering**: Date ranges use YYYY-MM-DD format. The `isWithinDateRange()` function handles ISO 8601 timestamps from Bitbucket.

## Development Commands

### Build and Run
```bash
go build -o server .
./server
```

### Run Directly
```bash
go run .
```

### Environment Setup
Copy `.env-example` to `.env` and fill in:
- BITBUCKET_EMAIL: Your Bitbucket email
- BITBUCKET_TOKEN: App password from Bitbucket (not your account password)
- BITBUCKET_WORKSPACE: Workspace slug/ID

### Test API Endpoints
```bash
# Get all commits across workspace
curl http://localhost:8080/commits

# Get commits for specific user with date filter
curl "http://localhost:8080/user-commits?account_id=<account_id>&start_date=2024-01-01&end_date=2024-12-31"

# Get workspace users
curl http://localhost:8080/repository-users

# Get user commit frequency
curl "http://localhost:8080/user-commit-frequency?account_id=<account_id>"
```

## API Endpoints

All endpoints support GET requests and return JSON:

- `GET /` - Serves dashboard.html
- `GET /commits?workspace={workspace}` - All commits for all repos in workspace
- `GET /commit?workspace={workspace}&repo={repo}&hash={hash}` - Specific commit
- `GET /pullrequests?workspace={workspace}&repo={repo}` - Merged PRs for a repo
- `GET /pullrequest?workspace={workspace}&repo={repo}&id={id}` - Specific PR
- `GET /repository-users?workspace={workspace}` - All workspace members
- `GET /user-commits?workspace={workspace}&account_id={account_id}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}` - User commits with optional date filter
- `GET /user-pullrequests?workspace={workspace}&account_id={account_id}` - User merged PRs
- `GET /user-commit-frequency?workspace={workspace}&account_id={account_id}` - Commit frequency statistics

The workspace parameter defaults to the configured BITBUCKET_WORKSPACE if omitted.

## Frontend (dashboard.html)

Single-page dashboard using:
- Tailwind CSS for styling
- Chart.js for visualizations (pie chart for commit distribution, bar chart for repo activity)
- Vanilla JavaScript for API calls and DOM manipulation

The dashboard hardcodes workspace 'eponce2710' in several places (lines 141, 711) - these should match your BITBUCKET_WORKSPACE when modifying.

## Important Notes

- **Rate Limiting**: The Bitbucket API has rate limits. When iterating over many repos, requests are sequential and can be slow.
- **Authentication**: Uses Basic Auth with email + app password (not OAuth).
- **Error Handling**: Errors from individual repos are logged but don't fail the entire request - allows partial results.
- **Pagination**: Current implementation does NOT handle paginated results from Bitbucket API - only first page of results is returned.
- **Binary**: There's a compiled binary `./server` that can be run directly without rebuilding.
