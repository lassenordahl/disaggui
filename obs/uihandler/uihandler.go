package uihandler

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

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

	// Create directory to store the downloaded bundle
	dir := "obsbundle"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// If stuff already exists in the obsbundle, remove it.
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}

	// Create a temporary file to store the downloaded zip file
	tempFile, err := os.Create(filepath.Join(dir, "bundle.zip"))
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	// Unzip the downloaded file
	if err := unzip(tempFile.Name(), dir); err != nil {
		return err
	}

	// Remove the temporary zip file
	if err := os.Remove(tempFile.Name()); err != nil {
		return err
	}

	// Move all the files in the `obsbundle/{version}` directory directly
	// into the `obsbundle` path.
	files, err = os.ReadDir(filepath.Join(dir, version))
	if err != nil {
		return err
	}

	for _, file := range files {
		src := filepath.Join(dir, version, file.Name())
		dest := filepath.Join(dir, file.Name())
		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}

	return nil
}

// Unzip extracts a zip archive to a destination directory
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fPath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
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

	// Serve assets with proper MIME types.
	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("obsbundle/assets")))
	r.PathPrefix("/assets/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		assetsHandler.ServeHTTP(w, r)
	}))

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("obsbundle"))))

	log.Println("UI is being served on port :8080")
}
