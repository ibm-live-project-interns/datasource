package syslogsim

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SyslogRecord struct {
	Raw       string    `json:"raw"`
	Priority  int       `json:"priority"`
	Timestamp time.Time `json:"timestamp"`
}

var fileLock sync.Mutex

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
		_ = json.Unmarshal(data, &records)
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
