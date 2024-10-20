package githubclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HttpClient   *http.Client
	Token        string
	GitHubAPIURL string
}

func NewClient() *Client {
	return &Client{
		HttpClient:   &http.Client{Timeout: 10 * time.Second},
		Token:        "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",
		GitHubAPIURL: "https://api.github.com/repositories",
	}
}

// Make authenticated request to GitHub API
func (c *Client) FetchRepositories(limit int) ([]Repository, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?per_page=%d", c.GitHubAPIURL, limit), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repositories: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repositories: %w", err)
	}

	return repositories, nil
}

// Fetch license for a specific repository
func (c *Client) FetchLicense(repoFullName string) (*License, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/license", repoFullName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch license for %s: %s", repoFullName, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var license License
	if err := json.Unmarshal(body, &license); err != nil {
		return nil, fmt.Errorf("failed to unmarshal license: %w", err)
	}

	return &license, nil
}

// Fetch languages for a specific repository
func (c *Client) FetchLanguages(repoFullName string) (map[string]int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/languages", repoFullName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch languages for %s: %s", repoFullName, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var languages map[string]int
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal languages: %w", err)
	}

	return languages, nil
}
