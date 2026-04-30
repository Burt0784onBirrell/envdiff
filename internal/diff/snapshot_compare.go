package diff

import "fmt"

// SnapshotDiff holds the changes between two snapshots.
type SnapshotDiff struct {
	NewMissingLeft  []string
	NewMissingRight []string
	NewMismatched   []Mismatch
	ResolvedKeys    []string
}

// CompareSnapshots returns the delta between an old and new snapshot.
func CompareSnapshots(old, new Snapshot) SnapshotDiff {
	oldMissingLeft := toSet(old.Result.MissingInLeft)
	oldMissingRight := toSet(old.Result.MissingInRight)
	oldMismatched := mismatchKeySet(old.Result.Mismatched)

	var sd SnapshotDiff

	for _, k := range new.Result.MissingInLeft {
		if !oldMissingLeft[k] {
			sd.NewMissingLeft = append(sd.NewMissingLeft, k)
		}
	}

	for _, k := range new.Result.MissingInRight {
		if !oldMissingRight[k] {
			sd.NewMissingRight = append(sd.NewMissingRight, k)
		}
	}

	newMismatched := mismatchKeySet(new.Result.Mismatched)
	for _, m := range new.Result.Mismatched {
		if !oldMismatched[m.Key] {
			sd.NewMismatched = append(sd.NewMismatched, m)
		}
	}

	// Keys present in old issues but not in new are resolved.
	for k := range oldMissingLeft {
		if !toSet(new.Result.MissingInLeft)[k] {
			sd.ResolvedKeys = append(sd.ResolvedKeys, k)
		}
	}
	for k := range oldMissingRight {
		if !toSet(new.Result.MissingInRight)[k] {
			sd.ResolvedKeys = append(sd.ResolvedKeys, k)
		}
	}
	for k := range oldMismatched {
		if !newMismatched[k] {
			sd.ResolvedKeys = append(sd.ResolvedKeys, k)
		}
	}

	return sd
}

// SnapshotDiffSummary returns a human-readable summary of the snapshot diff.
func SnapshotDiffSummary(sd SnapshotDiff) string {
	return fmt.Sprintf(
		"+%d new missing-left, +%d new missing-right, +%d new mismatched, %d resolved",
		len(sd.NewMissingLeft),
		len(sd.NewMissingRight),
		len(sd.NewMismatched),
		len(sd.ResolvedKeys),
	)
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func mismatchKeySet(mismatches []Mismatch) map[string]bool {
	s := make(map[string]bool, len(mismatches))
	for _, m := range mismatches {
		s[m.Key] = true
	}
	return s
}
