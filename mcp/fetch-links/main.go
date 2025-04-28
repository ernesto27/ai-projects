package main

import (
	"context"
	"encoding/json"
	"errors"
	fetch "fetchlinks/internal"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"fetch-links",
		"1.0.1",
	)

	// Add tool
	tool := mcp.NewTool("590_get-hacker-news-links",
		mcp.WithDescription("Get Links from Hacker News"),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of stories to fetch (top, new, best, ask, show, job) default to top"),
		),
	)

	tool2 := mcp.NewTool("590_get-infobae-links",
		mcp.WithDescription("Get Links from Infobae"),
	)

	// Add tool handler
	s.AddTool(tool, HackerNewsHandler)
	s.AddTool(tool2, InfobaeHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func HackerNewsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	typeValue, ok := request.Params.Arguments["type"].(string)
	if !ok {
		typeValue = "top" // Default to top if type is not provided or not a string
	}

	response, err := fetch.GetHackerNewsLinks(fetch.StoryType(typeValue))
	if err != nil {
		return nil, errors.New("failed to fetch links")
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("failed to marshal response to JSON")
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func InfobaeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	response, err := fetch.GetInfobaeLinks()
	if err != nil {
		return nil, errors.New("failed to fetch links")
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("failed to marshal response to JSON")
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
