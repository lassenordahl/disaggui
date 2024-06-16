package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./fingerprints.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS fingerprints (id INTEGER PRIMARY KEY AUTOINCREMENT, input TEXT, timestamp TEXT)")
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

func storeFingerprint(db *sql.DB, input, timestamp string) error {
	_, err := db.Exec("INSERT INTO fingerprints (input, timestamp) VALUES (?, ?)", input, timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert fingerprint: %v", err)
	}

	// Ensure no more than 300 rows
	err = enforceMaxRows(db, 300)
	if err != nil {
		return fmt.Errorf("failed to enforce max rows: %v", err)
	}

	return nil
}

func enforceMaxRows(db *sql.DB, maxRows int) error {
	row := db.QueryRow("SELECT COUNT(*) FROM fingerprints")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count rows: %v", err)
	}

	if count > maxRows {
		_, err := db.Exec("DELETE FROM fingerprints WHERE id = (SELECT id FROM fingerprints ORDER BY timestamp ASC LIMIT 1)")
		if err != nil {
			return fmt.Errorf("failed to delete oldest row: %v", err)
		}
	}
	return nil
}

type Fingerprint struct {
	Input     string `json:"input"`
	Timestamp string `json:"timestamp"`
}

type FingerprintPage struct {
	Fingerprints []Fingerprint `json:"fingerprints"`
	CurrentPage  int           `json:"current_page"`
	TotalPages   int           `json:"total_pages"`
}

func queryFingerprints(db *sql.DB, page, limit int) (FingerprintPage, error) {
	offset := (page - 1) * limit

	// Count total rows
	var totalRows int
	row := db.QueryRow("SELECT COUNT(*) FROM fingerprints")
	err := row.Scan(&totalRows)
	if err != nil {
		return FingerprintPage{}, fmt.Errorf("failed to count rows: %v", err)
	}

	totalPages := (totalRows + limit - 1) / limit // Calculate total pages

	query := fmt.Sprintf("SELECT input, timestamp FROM fingerprints ORDER BY timestamp DESC LIMIT %d OFFSET %d", limit, offset)
	rows, err := db.Query(query)
	if err != nil {
		return FingerprintPage{}, fmt.Errorf("failed to query fingerprints: %v", err)
	}
	defer rows.Close()

	var fingerprints []Fingerprint
	for rows.Next() {
		var fp Fingerprint
		if err := rows.Scan(&fp.Input, &fp.Timestamp); err != nil {
			return FingerprintPage{}, fmt.Errorf("failed to scan row: %v", err)
		}
		fingerprints = append(fingerprints, fp)
	}

	return FingerprintPage{
		Fingerprints: fingerprints,
		CurrentPage:  page,
		TotalPages:   totalPages,
	}, nil
}

type IntervalCount struct {
	Timestamp string `json:"timestamp"`
	Count     int    `json:"count"`
}

func getIntervalCounts(db *sql.DB) ([]IntervalCount, error) {
	rows, err := db.Query("SELECT timestamp FROM fingerprints ORDER BY timestamp ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query timestamps: %v", err)
	}
	defer rows.Close()

	var counts []IntervalCount
	var timestamp string
	for rows.Next() {
		err := rows.Scan(&timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan timestamp: %v", err)
		}

		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %v", err)
		}

		intervalKey := t.Truncate(30 * time.Second).Format("2006-01-02 15:04:05")
		found := false
		for i, count := range counts {
			if count.Timestamp == intervalKey {
				counts[i].Count++
				found = true
				break
			}
		}
		if !found {
			counts = append(counts, IntervalCount{Timestamp: intervalKey, Count: 1})
		}
	}

	return counts, nil
}
