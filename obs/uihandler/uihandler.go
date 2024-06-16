package uihandler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/gorilla/mux"
)

// DownloadBundle downloads the bundle for the specified version from the bucket server
func DownloadBundle(version string) error {
	url := fmt.Sprintf("http://localhost:8081/versions/%s", version)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download bundle: %s", resp.Status)
	}

	dir := "webbundle"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(dir, "index.html")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// GetLatestVersion fetches the latest patch version for a given major release from the bucket server
func GetLatestVersion(majorRelease string) (string, error) {
	resp, err := http.Get("http://localhost:8081/versions")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get versions: %s", resp.Status)
	}

	var versions []string
	for {
		var version string
		_, err := fmt.Fscanf(resp.Body, "%s\n", &version)
		if err != nil {
			break
		}
		versions = append(versions, version)
	}

	// Filter and sort versions
	var matchingVersions []string
	re := regexp.MustCompile(`^` + majorRelease + `\.\d+$`)
	for _, version := range versions {
		if re.MatchString(version) {
			matchingVersions = append(matchingVersions, version)
		}
	}
	sort.Strings(matchingVersions)

	if len(matchingVersions) == 0 {
		return "", fmt.Errorf("no matching versions found for major release %s", majorRelease)
	}

	return matchingVersions[len(matchingVersions)-1], nil
}

// Serve serves the latest bundle for the specified major release version using the provided router
func Serve(majorRelease string, r *mux.Router) {
	latestVersion, err := GetLatestVersion(majorRelease)
	if err != nil {
		log.Fatalf("Failed to get latest version: %v", err)
	}

	err = DownloadBundle(latestVersion)
	if err != nil {
		log.Fatalf("Failed to download bundle: %v", err)
	}

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("webbundle"))))

	log.Println("UI is being served on port :8080")
}
