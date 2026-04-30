package diff

import (
	"fmt"
	"io"
)

// WriteSnapshotReport writes a human-readable diff between two snapshots.
func WriteSnapshotReport(w io.Writer, old, new Snapshot, opts ReportOptions) {
	sd := CompareSnapshots(old, new)

	fmt.Fprintf(w, "Snapshot comparison\n")
	fmt.Fprintf(w, "  Old: %s (captured %s)\n", old.LeftPath+" vs "+old.RightPath, old.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "  New: %s (captured %s)\n", new.LeftPath+" vs "+new.RightPath, new.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(w)

	if len(sd.NewMissingLeft) > 0 {
		fmt.Fprintf(w, "New keys missing in left (%d):\n", len(sd.NewMissingLeft))
		for _, k := range sd.NewMissingLeft {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(sd.NewMissingRight) > 0 {
		fmt.Fprintf(w, "New keys missing in right (%d):\n", len(sd.NewMissingRight))
		for _, k := range sd.NewMissingRight {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(sd.NewMismatched) > 0 {
		fmt.Fprintf(w, "New mismatched keys (%d):\n", len(sd.NewMismatched))
		for _, m := range sd.NewMismatched {
			lv := maskValue(m.LeftValue, m.Key, opts)
			rv := maskValue(m.RightValue, m.Key, opts)
			fmt.Fprintf(w, "  ~ %s: %q != %q\n", m.Key, lv, rv)
		}
	}

	if len(sd.ResolvedKeys) > 0 {
		fmt.Fprintf(w, "Resolved keys (%d):\n", len(sd.ResolvedKeys))
		for _, k := range sd.ResolvedKeys {
			fmt.Fprintf(w, "  ✓ %s\n", k)
		}
	}

	if len(sd.NewMissingLeft)+len(sd.NewMissingRight)+len(sd.NewMismatched)+len(sd.ResolvedKeys) == 0 {
		fmt.Fprintln(w, "No changes between snapshots.")
	} else {
		fmt.Fprintln(w)
		fmt.Fprintln(w, SnapshotDiffSummary(sd))
	}
}
