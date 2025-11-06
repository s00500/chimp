//go:build mage

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
)

// UpdateDatastar fetches the latest Datastar release from jsDelivr CDN
func UpdateDatastar() error {
	fmt.Println("🔍 Fetching latest Datastar release from GitHub...")

	// Query GitHub API for latest release tag
	resp, err := http.Get("https://api.github.com/repos/starfederation/datastar/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode release info: %w", err)
	}

	version := release.TagName
	fmt.Printf("📦 Found version: %s\n", version)

	// Ensure static directory exists
	staticDir := "static"
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static directory: %w", err)
	}

	// Download from jsDelivr CDN (GitHub)
	downloadURL := fmt.Sprintf("https://cdn.jsdelivr.net/gh/starfederation/datastar@%s/bundles/datastar.js", version)
	outputPath := filepath.Join(staticDir, "datastar.min.js")

	if err := downloadFile(downloadURL, outputPath); err != nil {
		return fmt.Errorf("failed to download Datastar: %w", err)
	}

	fmt.Printf("✅ Updated %s to version %s\n", outputPath, version)
	return nil
}

// UpdateBasecoat fetches the latest BaseCoat assets from CDN
func UpdateBasecoat() error {
	fmt.Println("🔍 Fetching latest BaseCoat version from jsDelivr...")

	// Query jsDelivr API for latest version
	resp, err := http.Get("https://data.jsdelivr.com/v1/packages/npm/basecoat-css")
	if err != nil {
		return fmt.Errorf("failed to fetch jsDelivr package info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jsDelivr API returned status %d", resp.StatusCode)
	}

	var pkgInfo struct {
		Tags struct {
			Latest string `json:"latest"`
		} `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pkgInfo); err != nil {
		return fmt.Errorf("failed to decode package info: %w", err)
	}

	version := pkgInfo.Tags.Latest
	fmt.Printf("📦 Found version: %s\n", version)

	// Ensure static directory exists
	staticDir := "static"
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static directory: %w", err)
	}

	// Download CSS file
	cssURL := fmt.Sprintf("https://cdn.jsdelivr.net/npm/basecoat-css@%s/dist/basecoat.cdn.min.css", version)
	if err := downloadFile(cssURL, filepath.Join(staticDir, "basecoat.cdn.min.css")); err != nil {
		return fmt.Errorf("failed to download CSS: %w", err)
	}
	fmt.Printf("✅ Updated static/basecoat.cdn.min.css to version %s\n", version)

	// Download JS file
	jsURL := fmt.Sprintf("https://cdn.jsdelivr.net/npm/basecoat-css@%s/dist/js/all.min.js", version)
	if err := downloadFile(jsURL, filepath.Join(staticDir, "all.min.js")); err != nil {
		return fmt.Errorf("failed to download JS: %w", err)
	}
	fmt.Printf("✅ Updated static/all.min.js to version %s\n", version)

	return nil
}

// UpdateBasecoatSource fetches the latest BaseCoat source CSS from GitHub
func UpdateBasecoatSource() error {
	fmt.Println("🔍 Fetching latest BaseCoat source from GitHub...")

	// Query GitHub API for latest release tag
	resp, err := http.Get("https://api.github.com/repos/hunvreus/basecoat/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode release info: %w", err)
	}

	version := release.TagName
	fmt.Printf("📦 Found version: %s\n", version)

	// Ensure templates directory exists
	templatesDir := filepath.Join("cli", "chimp", "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Download source CSS from GitHub
	downloadURL := fmt.Sprintf("https://raw.githubusercontent.com/hunvreus/basecoat/%s/src/css/basecoat.css", version)
	outputPath := filepath.Join(templatesDir, "basecoat.css")

	if err := downloadFile(downloadURL, outputPath); err != nil {
		return fmt.Errorf("failed to download BaseCoat source: %w", err)
	}

	fmt.Printf("✅ Updated %s to version %s\n", outputPath, version)
	return nil
}

// UpdateAssets updates both Datastar and BaseCoat assets
func UpdateAssets() error {
	mg.Deps(UpdateDatastar, UpdateBasecoat, UpdateBasecoatSource)
	fmt.Println("\n🎉 All assets updated successfully!")
	return nil
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d for %s", resp.StatusCode, url)
	}

	outFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath, err)
	}

	return nil
}
