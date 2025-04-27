package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Base API URL
const baseAPI = "https://hacker-news.firebaseio.com/v0"

// Story type API endpoints
const (
	TopStoriesAPI  = baseAPI + "/topstories.json"
	NewStoriesAPI  = baseAPI + "/newstories.json"
	BestStoriesAPI = baseAPI + "/beststories.json"
	AskStoriesAPI  = baseAPI + "/askstories.json"
	ShowStoriesAPI = baseAPI + "/showstories.json"
	JobStoriesAPI  = baseAPI + "/jobstories.json"
	UpdatesAPI     = baseAPI + "/updates.json"
	UserAPI        = baseAPI + "/user/%s.json"
	ItemAPI        = baseAPI + "/item/%d.json"
)

// StoryType represents the type of stories to fetch
type StoryType string

// Available story types
const (
	TopStories  StoryType = "top"
	NewStories  StoryType = "new"
	BestStories StoryType = "best"
	AskStories  StoryType = "ask"
	ShowStories StoryType = "show"
	JobStories  StoryType = "job"
	Updates     StoryType = "updates"
)

type HackerNewsItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func GetHackerNewsLinks(storyType StoryType) ([]HackerNewsItem, error) {
	var apiURL string
	switch storyType {
	case TopStories:
		apiURL = TopStoriesAPI
	case NewStories:
		apiURL = NewStoriesAPI
	case BestStories:
		apiURL = BestStoriesAPI
	case AskStories:
		apiURL = AskStoriesAPI
	case ShowStories:
		apiURL = ShowStoriesAPI
	case JobStories:
		apiURL = JobStoriesAPI
	case Updates:
		apiURL = UpdatesAPI
	default:
		apiURL = TopStoriesAPI // Default to top stories
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch stories")
	}

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return nil, err
	}

	items := []HackerNewsItem{}
	for i, id := range storyIDs {
		if i >= 50 {
			break
		}

		itemResp, err := http.Get(fmt.Sprintf(ItemAPI, id))
		if err != nil {
			continue
		}
		defer itemResp.Body.Close()

		if itemResp.StatusCode != http.StatusOK {
			continue
		}

		var item HackerNewsItem
		if err := json.NewDecoder(itemResp.Body).Decode(&item); err != nil {
			continue
		}

		if item.URL != "" {
			items = append(items, item)
		}
	}

	return items, nil
}

// GetUserInfo fetches information about a specific user by ID
func GetUserInfo(userID string) (interface{}, error) {
	resp, err := http.Get(fmt.Sprintf(UserAPI, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	var userInfo interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
