package diff

import (
	"testing"
	"time"
)

func makeSnap(result Result) Snapshot {
	return Snapshot{
		Timestamp: time.Now().UTC(),
		LeftPath:  "a",
		RightPath: "b",
		Result:    result,
	}
}

func TestCompareSnapshots_NoChanges(t *testing.T) {
	r := Result{
		MissingInLeft:  []string{"X"},
		MissingInRight: []string{"Y"},
		Mismatched:     []Mismatch{{Key: "Z", LeftValue: "1", RightValue: "2"}},
	}
	sd := CompareSnapshots(makeSnap(r), makeSnap(r))
	if len(sd.NewMissingLeft)+len(sd.NewMissingRight)+len(sd.NewMismatched) != 0 {
		t.Errorf("expected no new issues, got %+v", sd)
	}
}

func TestCompareSnapshots_NewIssues(t *testing.T) {
	old := makeSnap(Result{})
	new := makeSnap(Result{
		MissingInLeft:  []string{"A"},
		MissingInRight: []string{"B"},
		Mismatched:     []Mismatch{{Key: "C", LeftValue: "x", RightValue: "y"}},
	})

	sd := CompareSnapshots(old, new)

	if len(sd.NewMissingLeft) != 1 || sd.NewMissingLeft[0] != "A" {
		t.Errorf("expected NewMissingLeft [A], got %v", sd.NewMissingLeft)
	}
	if len(sd.NewMissingRight) != 1 || sd.NewMissingRight[0] != "B" {
		t.Errorf("expected NewMissingRight [B], got %v", sd.NewMissingRight)
	}
	if len(sd.NewMismatched) != 1 || sd.NewMismatched[0].Key != "C" {
		t.Errorf("expected NewMismatched [C], got %v", sd.NewMismatched)
	}
}

func TestCompareSnapshots_ResolvedKeys(t *testing.T) {
	old := makeSnap(Result{
		MissingInLeft: []string{"A"},
		Mismatched:    []Mismatch{{Key: "B", LeftValue: "1", RightValue: "2"}},
	})
	new := makeSnap(Result{})

	sd := CompareSnapshots(old, new)

	if len(sd.ResolvedKeys) < 2 {
		t.Errorf("expected at least 2 resolved keys, got %v", sd.ResolvedKeys)
	}
}

func TestSnapshotDiffSummary(t *testing.T) {
	sd := SnapshotDiff{
		NewMissingLeft:  []string{"A"},
		NewMissingRight: []string{"B", "C"},
		NewMismatched:   []Mismatch{{Key: "D"}},
		ResolvedKeys:    []string{"E", "F", "G"},
	}
	summary := SnapshotDiffSummary(sd)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	expected := "+1 new missing-left, +2 new missing-right, +1 new mismatched, 3 resolved"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
