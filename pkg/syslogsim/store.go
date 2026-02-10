package syslogsim

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SyslogRecord represents a single persisted syslog event.
type SyslogRecord struct {
	Raw       string    `json:"raw"`       // The full RFC 5424 formatted syslog message
	Priority  int       `json:"priority"`  // Calculated priority (facility * 8 + severity)
	Timestamp time.Time `json:"timestamp"` // When the record was persisted (UTC)
}

// fileLock guards concurrent access to the syslog persistence file.
// NOTE: This is a package-level mutex shared across all file paths.
// For multi-file persistence, consider a per-path lock map.
var fileLock sync.Mutex

// SaveSyslogToFile appends a syslog record to the JSON array at the given path.
//
// NOTE: Same O(N^2) pattern as snmptrap/store.go — reads all records,
// appends one, writes everything back. For production use, consider JSON Lines
// (one JSON object per line) with append-only writes for O(1) per insert.
func SaveSyslogToFile(path string, raw string, priority int) error {
	fileLock.Lock()
	defer fileLock.Unlock()

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var records []SyslogRecord

	// Read existing file if it exists
	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		if unmarshalErr := json.Unmarshal(data, &records); unmarshalErr != nil {
			return fmt.Errorf("corrupt syslog file %s: %w", path, unmarshalErr)
		}
	}

	// Append new record
	records = append(records, SyslogRecord{
		Raw:       raw,
		Priority:  priority,
		Timestamp: time.Now().UTC(),
	})

	// Write back to file
	out, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}
