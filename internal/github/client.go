package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const GITHUB_API_URL = "https://api.github.com/repos/digitalpalidictionary/dpd-db/releases/latest"

type ReleaseInfo struct {
	Version     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	AssetURL    string
	HTMLURL     string `json:"html_url"`
}

type GitHubClient struct {
	APIURL string
	Client *http.Client
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		APIURL: GITHUB_API_URL,
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (gc *GitHubClient) GetLatestRelease() (*ReleaseInfo, error) {
	req, err := http.NewRequest("GET", gc.APIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "DPD-Updater-Go/1.0")

	resp, err := gc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var data struct {
		TagName     string    `json:"tag_name"`
		Name        string    `json:"name"`
		Body        string    `json:"body"`
		PublishedAt time.Time `json:"published_at"`
		HTMLURL     string    `json:"html_url"`
		Assets      []struct {
			Name      string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var assetURL string
	for _, asset := range data.Assets {
		if strings.HasPrefix(asset.Name, "dpd-goldendict") && strings.HasSuffix(asset.Name, ".zip") {
			assetURL = asset.DownloadURL
			break
		}
	}

	return &ReleaseInfo{
		Version:     data.TagName,
		Name:        data.Name,
		Body:        data.Body,
		PublishedAt: data.PublishedAt,
		AssetURL:    assetURL,
		HTMLURL:     data.HTMLURL,
	}, nil
}

func CompareVersions(current, latest string) int {
	if current == "unknown" {
		return -1
	}

	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")

	// Simple string comparison for now, can be improved with a semver library if needed
	if c < l {
		return -1
	}
	if c > l {
		return 1
	}
	return 0
}
