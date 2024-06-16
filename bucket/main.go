package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// ListFiles lists all the versions in the bundles directory
func listFiles() []string {
	files, err := os.ReadDir("bundles")
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	var versions []string
	for _, file := range files {
		// Remove the `.zip` suffix from the filename.
		version := strings.TrimSuffix(file.Name(), ".zip")

		// Don't add `.DS_Store` to the list of versions.
		if version == ".DS_Store" {
			continue
		}

		versions = append(versions, version)
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

// HandleBundle serves serves the `{version}.zip` bundle from the `bundles` directory.
// Users can download the full bundled zip from this file.
func handleBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]
	file := fmt.Sprintf("bundles/%s.zip", version)

	http.ServeFile(w, r, file)
}

func main() {
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
