package diff

import (
	"fmt"
	"io"
	"strings"
)

// WriteReport writes a human-readable diff report to the provided writer.
// leftName and rightName are labels for the two env files being compared.
func WriteReport(w io.Writer, result Result, leftName, rightName string) {
	if !result.HasDifferences() {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	if len(result.MissingInRight) > 0 {
		fmt.Fprintf(w, "Keys in %s but missing in %s:\n", leftName, rightName)
		for _, key := range result.MissingInRight {
			fmt.Fprintf(w, "  - %s\n", key)
		}
	}

	if len(result.MissingInLeft) > 0 {
		fmt.Fprintf(w, "Keys in %s but missing in %s:\n", rightName, leftName)
		for _, key := range result.MissingInLeft {
			fmt.Fprintf(w, "  + %s\n", key)
		}
	}

	if len(result.Mismatched) > 0 {
		fmt.Fprintln(w, "Mismatched values:")
		for _, m := range result.Mismatched {
			fmt.Fprintf(w, "  ~ %s\n", m.Key)
			fmt.Fprintf(w, "      %s: %s\n", leftName, maskValue(m.LeftValue))
			fmt.Fprintf(w, "      %s: %s\n", rightName, maskValue(m.RightValue))
		}
	}
}

// maskValue replaces the middle portion of a value with asterisks to avoid
// leaking sensitive data in reports.
func maskValue(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
