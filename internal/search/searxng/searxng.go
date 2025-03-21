package searxng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client represents a SearxNG API client
type Client struct {
	baseURL string // The base URL of the SearxNG instance
}

// NewClient creates a new SearxNG client with the specified base URL
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

// SearchOptions represents the available search parameters for SearxNG
type SearchOptions struct {
	Categories []string `json:"categories,omitempty"` // List of categories to search in (e.g., "general", "images")
	Engines    []string `json:"engines,omitempty"`    // List of search engines to use
	Language   string   `json:"language,omitempty"`   // Language code for search results (e.g., "en-US")
	PageNo     int      `json:"pageno,omitempty"`     // Page number for paginated results (starts at 1)
}

// SearchResult represents a single result from a SearxNG search
type SearchResult struct {
	Title        string `json:"title"`                   // Title of the search result
	URL          string `json:"url"`                     // URL of the search result
	ImgSrc       string `json:"img_src,omitempty"`       // Source URL of the image (for image results)
	ThumbnailSrc string `json:"thumbnail_src,omitempty"` // Source URL of the thumbnail
	Thumbnail    string `json:"thumbnail,omitempty"`     // Alternative thumbnail URL
	Content      string `json:"content,omitempty"`       // Snippet or description of the result
	Author       string `json:"author,omitempty"`        // Author of the content (if available)
	IframeSrc    string `json:"iframe_src,omitempty"`    // Source URL for iframe content
}

// SearchResponse represents the complete response from SearxNG including results and suggestions
type SearchResponse struct {
	Results     []*SearchResult `json:"results"`               // List of search results
	Suggestions []string        `json:"suggestions,omitempty"` // Search suggestions based on the query
}

// Search performs a search query against the SearxNG instance
//
// Parameters:
//   - query: The search term or phrase
//   - opts: Optional search parameters (can be nil for default settings)
//
// Returns:
//   - *SearchResponse: Contains search results and suggestions
//   - error: Any error that occurred during the search
func (c *Client) Search(query string, opts *SearchOptions) (*SearchResponse, error) {
	// Construct the base URL with the search query
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Create query parameters
	params := url.Values{}
	params.Set("format", "json")
	params.Set("q", query)

	// Add optional parameters if provided
	if opts != nil {
		if len(opts.Categories) > 0 {
			params.Set("categories", strings.Join(opts.Categories, ","))
		}

		if len(opts.Engines) > 0 {
			params.Set("engines", strings.Join(opts.Engines, ","))
		}

		if opts.Language != "" {
			params.Set("language", opts.Language)
		}

		if opts.PageNo > 0 {
			params.Set("pageno", fmt.Sprintf("%d", opts.PageNo))
		}
	}

	// Construct the final URL
	searchURL := baseURL.JoinPath("search")
	searchURL.RawQuery = params.Encode()

	// Create and execute the HTTP request
	req, err := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}
