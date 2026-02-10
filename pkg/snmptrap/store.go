package snmptrap

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// fileLock guards concurrent access to the trap persistence file.
// NOTE: This is a package-level mutex, which means all callers share a single
// lock regardless of which file path they target. For multi-file persistence,
// consider a per-path lock map.
var fileLock sync.Mutex

// SaveTrapToFile appends a trap to the JSON array stored at the given path.
//
// NOTE: The current implementation reads the entire file, deserializes all
// records, appends one, and writes everything back. This is O(N) per append
// and results in O(N^2) total cost over N inserts. For production workloads,
// consider using JSON Lines (one JSON object per line) with append-only writes,
// which would make each insert O(1).
func SaveTrapToFile(path string, trap Trap) error {
	fileLock.Lock()
	defer fileLock.Unlock()

	var traps []Trap

	// Read existing file
	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		if unmarshalErr := json.Unmarshal(data, &traps); unmarshalErr != nil {
			return fmt.Errorf("corrupt trap file %s: %w", path, unmarshalErr)
		}
	}

	// Append new trap
	traps = append(traps, trap)

	// Write back
	out, err := json.MarshalIndent(traps, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}
