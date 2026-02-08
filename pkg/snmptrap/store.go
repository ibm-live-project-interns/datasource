package snmptrap

import (
	"encoding/json"
	"os"
	"sync"
)

var fileLock sync.Mutex

func SaveTrapToFile(path string, trap Trap) error {
	fileLock.Lock()
	defer fileLock.Unlock()

	var traps []Trap

	// Read existing file
	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &traps)
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
