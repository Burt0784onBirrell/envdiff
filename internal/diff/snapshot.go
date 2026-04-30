package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved diff result at a point in time.
type Snapshot struct {
	Timestamp   time.Time  `json:"timestamp"`
	LeftPath    string     `json:"left_path"`
	RightPath   string     `json:"right_path"`
	Result      Result     `json:"result"`
}

// SaveSnapshot writes the current diff result to a JSON snapshot file.
func SaveSnapshot(path, leftPath, rightPath string, result Result) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		LeftPath:  leftPath,
		RightPath: rightPath,
		Result:    result,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// LoadSnapshot reads a snapshot file and returns the stored Snapshot.
func LoadSnapshot(path string) (Snapshot, error) {
	var snap Snapshot

	data, err := os.ReadFile(path)
	if err != nil {
		return snap, fmt.Errorf("snapshot: read failed: %w", err)
	}

	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return snap, nil
}
