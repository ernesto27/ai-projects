# AGENTS.md - Bitbucket Analytics API & Dashboard

## Change Review Requirements
- **ALWAYS** show proposed changes for review before making any modifications
- Display git-style diffs or full file content for user approval
- Wait for explicit permission ("yes", "proceed", etc.) before creating/modifying files
- Never make changes without user confirmation

## Build/Lint/Test Commands
- **Build**: `go build -o server .`
- **Run**: `go run .`
- **Run with Docker**: `docker-compose up --build`
- **Watch mode**: `docker-compose watch`
- **Tests**: No test files present - run manually or add Go tests as needed

## Architecture & Structure
- **Backend**: Go 1.25.1 web service with flat architecture
- **main.go**: HTTP server setup, route registration, config loading
- **bitbucket.go**: Bitbucket API client with authentication and data fetching
- **controllers.go**: HTTP handlers for all endpoints
- **types.go**: Data structures for API requests/responses
- **frontend/**: Single-page dashboard (HTML/JS with Chart.js)
- **Data flow**: HTTP handlers → API client → Bitbucket API → JSON responses

## Code Style Guidelines
- **Formatting**: Standard Go formatting (use `gofmt`)
- **Imports**: Group standard library first, then third-party
- **Naming**: PascalCase for exported types/functions, camelCase for unexported
- **Error handling**: Log errors but continue processing for partial results
- **JSON**: Use struct tags for all API types, nested anonymous structs for complex responses
- **Functions**: Descriptive names like `handleGetCommits`, `getRepositories`
- **Comments**: Document exported types and functions

## Environment Setup
- Copy `.env-example` to `.env` with BITBUCKET_EMAIL, BITBUCKET_TOKEN, BITBUCKET_WORKSPACE
- Server runs on port 8080 (configurable via PORT env var)

## CLAUDE.md Integration
This file incorporates guidance from CLAUDE.md including:
- Project overview as Bitbucket Analytics API & Dashboard
- Detailed architecture patterns and data flow
- API endpoints documentation
- Development commands and environment setup
- Important notes on rate limiting, authentication, and pagination
