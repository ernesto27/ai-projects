package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// LinkItem represents a link with a title and URL
type LinkItem struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source,omitempty"` // Optional field to identify the source
}

func GetHackerNewsLinks(storyType StoryType) ([]LinkItem, error) {
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

	items := []LinkItem{}
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

		var hackerNewsItem LinkItem
		if err := json.NewDecoder(itemResp.Body).Decode(&hackerNewsItem); err != nil {
			continue
		}

		if hackerNewsItem.URL != "" {
			items = append(items, LinkItem{
				Title:  hackerNewsItem.Title,
				URL:    hackerNewsItem.URL,
				Source: "hackernews",
			})
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

// GetInfobaeLinks scrapes news links from Infobae's website
func GetInfobaeLinks() ([]LinkItem, error) {
	const infobaeURL = "https://www.infobae.com"

	client := &http.Client{}
	req, err := http.NewRequest("GET", infobaeURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Infobae: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Infobae with status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Infobae response body: %w", err)
	}

	bodyReader := bytes.NewReader(bodyBytes)

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Infobae HTML: %w", err)
	}

	var newsItems []LinkItem

	selectors := []string{
		"a.headline",
		"article a",
		".story-card a",
		".headline a",
		"h1 a",
		"h2 a",
		"h3 a",
		".article-title a",
		"a[data-track-name]",
		".d23-story-card a",
		".headline-link",
		".card-headline a",
		"a.card-title",
		".entry-box-titles a",
		"a.title",
		"a[href*='/america/']",
		"a[href*='/economia/']",
		"a[href*='/politica/']",
		"a[href*='/tendencias/']",
	}

	selectorQuery := strings.Join(selectors, ", ")

	doc.Find(selectorQuery).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url, exists := s.Attr("href")

		title = strings.TrimSpace(title)

		// Only add items with non-empty title and URL
		if exists && title != "" && url != "" {
			if !strings.HasPrefix(url, "http") {
				url = infobaeURL + url
			}

			// Check if this item is already in our list (avoid duplicates)
			isDuplicate := false
			for _, item := range newsItems {
				if item.URL == url {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				newsItems = append(newsItems, LinkItem{
					Title:  title,
					URL:    url,
					Source: "infobae",
				})
			}
		}
	})

	// Debug information
	fmt.Printf("Found %d links from Infobae\n", len(newsItems))
	fmt.Println(newsItems)

	// Limit to 50 items if there are more
	if len(newsItems) > 100 {
		newsItems = newsItems[:100]
	}

	return newsItems, nil
}
