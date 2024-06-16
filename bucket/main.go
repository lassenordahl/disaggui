package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// GenerateFakeBundle creates a fake index.html file for a given version
func generateFakeBundle(version string) error {
	dir := filepath.Join("bundles", version)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(dir, "index.html")
	content := fmt.Sprintf("<html><body><h1>Version %s</h1></body></html>", version)
	return os.WriteFile(filePath, []byte(content), 0644)
}

// ListFiles lists all the versions in the bundles directory
func listFiles() []string {
	files, err := os.ReadDir("bundles")
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	var versions []string
	for _, file := range files {
		if file.IsDir() {
			versions = append(versions, file.Name())
		}
	}
	return versions
}

// HandleVersionsList handles the /versions endpoint to list all available versions
func handleVersionsList(w http.ResponseWriter, r *http.Request) {
	versions := listFiles()
	for _, version := range versions {
		fmt.Fprintln(w, version)
	}
}

// HandleBundle serves versioned bundles from a local directory
func handleBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]
	filePath := filepath.Join("bundles", version, "index.html")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func main() {
	// Generate some fake bundles
	versions := []string{"v1.0.0", "v1.0.1", "v1.1.0", "v1.2.3"}
	for _, version := range versions {
		if err := generateFakeBundle(version); err != nil {
			log.Fatalf("Failed to generate bundle for version %s: %v", version, err)
		}
	}

	// Create a new router
	r := mux.NewRouter()

	// Define the routes
	r.HandleFunc("/versions", handleVersionsList).Methods("GET")
	r.HandleFunc("/versions/{version}", handleBundle).Methods("GET")

	// Start the server
	log.Println("Serving at http://localhost:8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
