package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *server) listFingerprints(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}

	fingerprintPage, err := queryFingerprints(s.db, page, limit)
	if err != nil {
		http.Error(w, "Failed to query fingerprints", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(fingerprintPage)
}

func (s *server) listFingerprintCounts(w http.ResponseWriter, r *http.Request) {
	counts, err := getIntervalCounts(s.db)
	if err != nil {
		println(err.Error())
		http.Error(w, "Failed to count statements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}
