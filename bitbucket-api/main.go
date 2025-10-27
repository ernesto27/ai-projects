package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var config Config

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	config = Config{
		Email:     os.Getenv("BITBUCKET_EMAIL"),
		Token:     os.Getenv("BITBUCKET_TOKEN"),
		Workspace: os.Getenv("BITBUCKET_WORKSPACE"),
	}

	if config.Email == "" || config.Token == "" || config.Workspace == "" {
		log.Fatal("Please set BITBUCKET_EMAIL, BITBUCKET_TOKEN, and BITBUCKET_WORKSPACE environment variables")
	}

	// Create HTTP router
	mux := http.NewServeMux()

	// Serve static files from frontend directory
	fs := http.FileServer(http.Dir("frontend"))
	mux.Handle("/dashboard.js", fs)

	// Register endpoints
	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/commits", handleGetCommits)
	mux.HandleFunc("/commit", handleGetCommit)
	mux.HandleFunc("/repository-commit", handleRepositoryCommit)
	mux.HandleFunc("/pullrequests", handleGetPullRequests)
	mux.HandleFunc("/pullrequest", handleGetPullRequest)
	mux.HandleFunc("/repository-users", handleGetRepositoryUsers)
	mux.HandleFunc("/user-commits", handleGetUserCommits)
	mux.HandleFunc("/user-pullrequests", handleGetUserPullRequests)
	mux.HandleFunc("/user-commit-frequency", handleGetUserCommitFrequency)
	mux.HandleFunc("/languages", handleGetLanguages)

	// Start server
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Starting server on port %s...", port)
	log.Printf("Endpoints:")
	log.Printf("  GET / - Dashboard")
	log.Printf("  GET /commits?workspace={workspace} - Get commits for all repos in workspace")
	log.Printf("  GET /commit?workspace={workspace}&repo={repo}&hash={hash} - Get specific commit")
	log.Printf("  GET /repository-commit?workspace={workspace}&repo={repo} - Get all commits for a repository")
	log.Printf("  GET /pullrequests?workspace={workspace}&repo={repo} - Get pull requests for a repo")
	log.Printf("  GET /pullrequest?workspace={workspace}&repo={repo}&id={id} - Get specific pull request")
	log.Printf("  GET /repository-users?workspace={workspace} - Get all users in workspace")
	log.Printf("  GET /user-commits?workspace={workspace}&account_id={account_id}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD} - Get commits by user with date filter")
	log.Printf("  GET /user-pullrequests?workspace={workspace}&account_id={account_id} - Get merged pull requests by user account ID")
	log.Printf("  GET /user-commit-frequency?workspace={workspace}&account_id={account_id} - Get commit frequency statistics for all history")
	log.Printf("  GET /languages?workspace={workspace}&repo={repo} - Get language statistics for repositories")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
