# Bitbucket Analytics API & Dashboard

A Go web service that fetches repository data from the Bitbucket API and presents analytics through both a REST API and an interactive HTML dashboard.

## Features

- Interactive dashboard with real-time analytics
- User commit tracking and statistics
- Pull request analytics (merged and open)
- Commit frequency analysis
- Programming language statistics per repository
- Beautiful dark-themed UI with charts

## Requirements

### For Go Development
- **Go**: 1.25.1 or higher
- **Git**: For version control

### For Docker Development
- **Docker**: 20.10 or higher
- **Docker Compose**: v2.22 or higher (for watch feature)

## Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd bitbucket-api
```

### 2. Configure Environment Variables

Create a `.env` file from the example:

```bash
cp .env-example .env
```

Edit `.env` and add your Bitbucket credentials:

```env
BITBUCKET_EMAIL=your-email@example.com
BITBUCKET_TOKEN=your-app-password
BITBUCKET_WORKSPACE=your-workspace-slug
```

**Note**: The `BITBUCKET_TOKEN` should be a Bitbucket App Password, not your account password.

#### Creating a Bitbucket App Password:
1. Go to Bitbucket → Settings → App passwords
2. Create a new app password with permissions:
   - **Repositories**: Read
   - **Pull requests**: Read
   - **Workspace membership**: Read

## Development

### Option 1: Develop with Go (Native)


#### Run the Application

```bash
go run .
```


The server will start on `http://localhost:8080`


### Option 2: Develop with Docker Compose

#### Build and Run

```bash
docker-compose up --build
```

#### Run with Watch Mode (Auto-reload)

Docker Compose v2.22+ supports automatic rebuild on file changes:

```bash
docker-compose watch
```

This will:
- **Rebuild** the container when `.go` files or `.env` change
- **Sync** frontend files (HTML/JS/CSS) without rebuilding for faster updates

#### Stop the Application

```bash
docker-compose down
```

#### View Logs

```bash
docker-compose logs -f
```

#### Rebuild from Scratch

```bash
docker-compose down
docker-compose build --no-cache
docker-compose up
```

## API Endpoints

### Repository Data
- `GET /` - Dashboard UI
- `GET /commits?workspace={workspace}` - All commits for all repos
- `GET /commit?workspace={workspace}&repo={repo}&hash={hash}` - Specific commit
- `GET /pullrequests?workspace={workspace}&repo={repo}` - Pull requests for a repo
- `GET /pullrequest?workspace={workspace}&repo={repo}&id={id}` - Specific pull request

### User Analytics
- `GET /repository-users?workspace={workspace}` - All workspace members
- `GET /user-commits?workspace={workspace}&account_id={account_id}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}` - User commits with optional date filter
- `GET /user-pullrequests?workspace={workspace}&account_id={account_id}` - User merged pull requests
- `GET /user-commit-frequency?workspace={workspace}&account_id={account_id}` - Commit frequency statistics

### Language Statistics
- `GET /languages?workspace={workspace}&repo={repo}` - Language statistics for repositories

**Note**: The `workspace` parameter defaults to the configured `BITBUCKET_WORKSPACE` if omitted.

## Testing API Endpoints

```bash
# Get all commits across workspace
curl http://localhost:8080/commits

# Get commits for specific user with date filter
curl "http://localhost:8080/user-commits?account_id=<account_id>&start_date=2024-01-01&end_date=2024-12-31"

# Get workspace users
curl http://localhost:8080/repository-users

# Get user commit frequency
curl "http://localhost:8080/user-commit-frequency?account_id=<account_id>"

# Get language statistics for a repository
curl "http://localhost:8080/languages?repo=cicd-demo"
```

## Dashboard Features

Access the dashboard at `http://localhost:8080`

The dashboard provides:
- **Overview Statistics**: Total repos, commits, contributors
- **User Analytics**: Filter by workspace member and time period
  - Commits per repository
  - Merged and open pull requests
  - Commit frequency statistics
- **Repository Analytics**: Select individual repositories to see:
  - Contributor breakdown
  - Commit distribution chart
  - Language statistics chart
- **Repository Activity**: Bar chart showing commits per repository

## Technology Stack

### Backend
- **Go 1.25.1**: Main programming language
- **Standard Library**: HTTP server and JSON handling
- **godotenv**: Environment variable management

### Frontend
- **Tailwind CSS**: Styling framework
- **Chart.js**: Data visualization
- **Vanilla JavaScript**: Dashboard interactivity

### DevOps
- **Docker**: Containerization
- **Docker Compose**: Container orchestration with watch mode





