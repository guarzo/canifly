package skillplans

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

// GitHubDownloader handles downloading skill plans from a GitHub repository
type GitHubDownloader struct {
	repoURL        string
	httpClient     *http.Client
	logger         interfaces.Logger
	circuitBreaker *CircuitBreaker
}

// NewGitHubDownloader creates a new GitHub downloader service
func NewGitHubDownloader(repoURL string, logger interfaces.Logger) *GitHubDownloader {
	return &GitHubDownloader{
		repoURL: repoURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:         logger,
		circuitBreaker: NewCircuitBreaker(3, 5*time.Minute), // Open after 3 failures, reset after 5 minutes
	}
}

// DownloadPlans downloads all skill plans from GitHub to the specified directory
func (g *GitHubDownloader) DownloadPlans(destDir string) error {
	if g.repoURL == "" {
		return fmt.Errorf("GitHub repository URL not configured")
	}

	g.logger.Infof("Downloading skill plans from GitHub to %s", destDir)

	// Fetch the directory listing from GitHub API
	// Convert raw GitHub URL to API URL
	apiURL := g.convertToAPIURL()
	if apiURL == "" {
		// Fallback to predefined list if API URL conversion fails
		return g.downloadPredefinedPlans(destDir)
	}

	// Fetch the list of files from GitHub API with circuit breaker protection
	var resp *http.Response
	err := g.circuitBreaker.Call(func() error {
		var err error
		resp, err = g.httpClient.Get(apiURL)
		return err
	})

	if err != nil {
		g.logger.Warnf("Failed to fetch directory listing from API, falling back to predefined list: %v", err)
		return g.downloadPredefinedPlans(destDir)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.logger.Warnf("GitHub API returned %d, falling back to predefined list", resp.StatusCode)
		return g.downloadPredefinedPlans(destDir)
	}

	// Parse the GitHub API response to get file names
	var files []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&files); err != nil {
		g.logger.Warnf("Failed to parse GitHub API response, falling back to predefined list: %v", err)
		return g.downloadPredefinedPlans(destDir)
	}

	// Download each .txt file
	downloadCount := 0
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(file.Name, ".txt") {
			planName := strings.TrimSuffix(file.Name, ".txt")
			destPath := filepath.Join(destDir, file.Name)

			// Skip if file already exists and is recent (< 24 hours old)
			if info, err := os.Stat(destPath); err == nil {
				if time.Since(info.ModTime()) < 24*time.Hour {
					g.logger.Debugf("Skipping %s - file is recent", planName)
					continue
				}
			}

			if err := g.DownloadPlan(planName, destPath); err != nil {
				if err == ErrCircuitBreakerOpen {
					g.logger.Warnf("Circuit breaker is open, stopping further downloads")
					break
				}
				g.logger.Warnf("Failed to download %s: %v", planName, err)
				// Continue with other plans
			} else {
				g.logger.Debugf("Downloaded %s successfully", planName)
				downloadCount++
			}
		}
	}

	if downloadCount == 0 {
		g.logger.Warn("No skill plans were downloaded, falling back to predefined list")
		return g.downloadPredefinedPlans(destDir)
	}

	return nil
}

// convertToAPIURL converts a raw GitHub URL to the GitHub API URL
func (g *GitHubDownloader) convertToAPIURL() string {
	// Convert URLs like https://raw.githubusercontent.com/user/repo/branch/path
	// to https://api.github.com/repos/user/repo/contents/path?ref=branch
	if strings.Contains(g.repoURL, "raw.githubusercontent.com") {
		parts := strings.Split(g.repoURL, "/")
		if len(parts) >= 7 {
			user := parts[3]
			repo := parts[4]
			branch := parts[5]
			path := strings.Join(parts[6:], "/") + "/plans"
			return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", user, repo, path, branch)
		}
	}
	return ""
}

// downloadPredefinedPlans falls back to downloading a predefined list of plans
func (g *GitHubDownloader) downloadPredefinedPlans(destDir string) error {
	planNames := []string{
		"Bifrost", "Flycatcher", "Jaguar", "Keres",
		"Kiki", "Leshak", "Magic_14", "Manticore", "Naga",
	}

	for _, name := range planNames {
		destPath := filepath.Join(destDir, name+".txt")

		// Skip if file already exists and is recent (< 24 hours old)
		if info, err := os.Stat(destPath); err == nil {
			if time.Since(info.ModTime()) < 24*time.Hour {
				g.logger.Debugf("Skipping %s - file is recent", name)
				continue
			}
		}

		if err := g.DownloadPlan(name, destPath); err != nil {
			if err == ErrCircuitBreakerOpen {
				g.logger.Warnf("Circuit breaker is open, stopping further downloads")
				break
			}
			g.logger.Warnf("Failed to download %s: %v", name, err)
			// Continue with other plans
		} else {
			g.logger.Debugf("Downloaded %s successfully", name)
		}
	}

	return nil
}

// DownloadPlan downloads a single skill plan from GitHub
func (g *GitHubDownloader) DownloadPlan(planName, destPath string) error {
	if g.repoURL == "" {
		return fmt.Errorf("GitHub repository URL not configured")
	}

	// Use circuit breaker to protect against repeated failures
	return g.circuitBreaker.Call(func() error {
		// URL-escape the plan name to handle special characters
		escapedPlanName := url.QueryEscape(planName)
		downloadURL := fmt.Sprintf("%s/plans/%s.txt", g.repoURL, escapedPlanName)

		resp, err := g.httpClient.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("failed to fetch plan: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}

		// Ensure the destination directory exists
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Create the file
		file, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

		// Copy the response body to the file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	})
}
