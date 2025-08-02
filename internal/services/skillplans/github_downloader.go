package skillplans

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

// GitHubDownloader handles downloading skill plans from a GitHub repository
type GitHubDownloader struct {
	repoURL    string
	httpClient *http.Client
	logger     interfaces.Logger
}

// NewGitHubDownloader creates a new GitHub downloader service
func NewGitHubDownloader(repoURL string, logger interfaces.Logger) *GitHubDownloader {
	return &GitHubDownloader{
		repoURL: repoURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// DownloadPlans downloads all skill plans from GitHub to the specified directory
func (g *GitHubDownloader) DownloadPlans(destDir string) error {
	planNames := []string{
		"Bifrost", "Flycatcher", "Jaguar", "Keres",
		"Kiki", "Leshak", "Magic_14", "Manticore", "Naga",
	}

	g.logger.Infof("Downloading skill plans from GitHub to %s", destDir)

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

	url := fmt.Sprintf("%s/plans/%s.txt", g.repoURL, planName)

	resp, err := g.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch plan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
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
}