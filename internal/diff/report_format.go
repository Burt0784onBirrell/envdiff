package diff

import (
	"fmt"
	"io"
	"strings"
)

// WriteFormattedReport writes the diff result using the specified OutputFormat.
func WriteFormattedReport(w io.Writer, result *Result, opts ReportOptions, format OutputFormat) error {
	if format == FormatMarkdown {
		fmt.Fprint(w, MarkdownTableHeader())
	}

	for _, key := range result.MissingInRight {
		line := FormatMissingLine(format, key, opts.RightLabel)
		fmt.Fprintln(w, line)
	}

	for _, key := range result.MissingInLeft {
		line := FormatMissingLine(format, key, opts.LeftLabel)
		fmt.Fprintln(w, line)
	}

	for _, mm := range result.Mismatched {
		leftVal := mm.LeftValue
		rightVal := mm.RightValue
		if opts.MaskSecrets && IsSensitiveKey(mm.Key) {
			leftVal = maskValue(leftVal)
			rightVal = maskValue(rightVal)
		}
		line := FormatMismatchLine(format, mm.Key, leftVal, rightVal)
		fmt.Fprintln(w, line)
	}

	if format == FormatText {
		total := len(result.MissingInRight) + len(result.MissingInLeft) + len(result.Mismatched)
		if total == 0 {
			fmt.Fprintln(w, "No differences found.")
		} else {
			fmt.Fprintf(w, "\nTotal differences: %d\n", total)
		}
	}

	return nil
}

// SummaryLine returns a single-line human-readable summary of the diff result.
func SummaryLine(result *Result, leftLabel, rightLabel string) string {
	parts := []string{}
	if n := len(result.MissingInRight); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in %s", n, rightLabel))
	}
	if n := len(result.MissingInLeft); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in %s", n, leftLabel))
	}
	if n := len(result.Mismatched); n > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", n))
	}
	if len(parts) == 0 {
		return "No differences found."
	}
	return strings.Join(parts, ", ")
}
