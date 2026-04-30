package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baseSnapshotResult() Result {
	return Result{
		MissingInLeft:  []string{"A"},
		MissingInRight: []string{"B"},
		Mismatched:     []Mismatch{{Key: "C", LeftValue: "x", RightValue: "y"}},
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	result := baseSnapshotResult()
	err := SaveSnapshot(path, ".env.left", ".env.right", result)
	if err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if snap.LeftPath != ".env.left" {
		t.Errorf("expected left path .env.left, got %s", snap.LeftPath)
	}
	if snap.RightPath != ".env.right" {
		t.Errorf("expected right path .env.right, got %s", snap.RightPath)
	}
	if len(snap.Result.MissingInLeft) != 1 || snap.Result.MissingInLeft[0] != "A" {
		t.Errorf("unexpected MissingInLeft: %v", snap.Result.MissingInLeft)
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadSnapshot_InvalidPath(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSnapshot_BadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for bad JSON")
	}
}

func TestSaveSnapshot_TimestampIsUTC(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	before := time.Now().UTC()
	SaveSnapshot(path, "a", "b", Result{})
	after := time.Now().UTC()

	snap, _ := LoadSnapshot(path)
	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", snap.Timestamp, before, after)
	}
}
