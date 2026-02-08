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
			Name        string `json:"name"`
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

	// Try to parse both as dates first
	currentDate, err1 := parseDate(current)
	latestDate, err2 := parseDate(latest)

	if err1 == nil && err2 == nil {
		// Both are dates, compare date-only (ignore time component)
		currentY, currentM, currentD := currentDate.Date()
		latestY, latestM, latestD := latestDate.Date()

		if currentY < latestY || (currentY == latestY && currentM < latestM) || (currentY == latestY && currentM == latestM && currentD < latestD) {
			return -1
		}
		if currentY > latestY || (currentY == latestY && currentM > latestM) || (currentY == latestY && currentM == latestM && currentD > latestD) {
			return 1
		}
		return 0
	}

	// If current is a date but latest is not, try semver format
	if err1 == nil && err2 != nil {
		// Parse GitHub version (semver with date as patch): "v0.3.20260202"
		latest = strings.TrimPrefix(latest, "v")
		parts := strings.Split(latest, ".")
		if len(parts) >= 3 {
			// Patch is the 3rd part and is YYYYMMDD format
			patch := parts[2]
			if len(patch) >= 8 {
				// Extract first 8 digits as YYYYMMDD
				dateStr := patch[:8]
				if latestDate, err := time.Parse("20060102", dateStr); err == nil {
					// Compare dates (date only, ignore time)
					currentY, currentM, currentD := currentDate.Date()
					latestY, latestM, latestD := latestDate.Date()

					if currentY < latestY || (currentY == latestY && currentM < latestM) || (currentY == latestY && currentM == latestM && currentD < latestD) {
						return -1
					}
					if currentY > latestY || (currentY == latestY && currentM > latestM) || (currentY == latestY && currentM == latestM && currentD > latestD) {
						return 1
					}
					return 0
				}
			}
		}
	}

	// Fallback to string comparison
	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")
	if c < l {
		return -1
	}
	if c > l {
		return 1
	}
	return 0
}

// parseDate tries multiple date formats and returns the parsed time
func parseDate(s string) (time.Time, error) {
	// Common date formats to try
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02",
		"2006.01.02",
		"2006/01/02",
		"20060102",
	}

	// Try parsing the full string first
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	// Try to extract date from the end of version strings like "0.3.20260202"
	// Look for 8 consecutive digits at the end (YYYYMMDD format)
	for i := len(s) - 8; i >= 0; i-- {
		if i+8 <= len(s) {
			candidate := s[i : i+8]
			// Check if it's all digits
			if isAllDigits(candidate) {
				if t, err := time.Parse("20060102", candidate); err == nil {
					return t, nil
				}
			}
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// isAllDigits checks if a string contains only digits
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
