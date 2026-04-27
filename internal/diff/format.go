package diff

import (
	"fmt"
	"strings"
)

// OutputFormat defines the format for report output.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatMarkdown OutputFormat = "markdown"
)

// ParseOutputFormat parses a string into an OutputFormat.
// Returns FormatText and an error if the format is unrecognized.
func ParseOutputFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	default:
		return FormatText, fmt.Errorf("unknown output format %q: must be one of text, json, markdown", s)
	}
}

// FormatMissingLine returns a formatted string for a missing key entry.
func FormatMissingLine(format OutputFormat, key, side string) string {
	switch format {
	case FormatJSON:
		return fmt.Sprintf(`{"type":"missing","key":%q,"missing_in":%q}`, key, side)
	case FormatMarkdown:
		return fmt.Sprintf("| MISSING | `%s` | — | %s |", key, side)
	default:
		return fmt.Sprintf("MISSING  %-30s (not in %s)", key, side)
	}
}

// FormatMismatchLine returns a formatted string for a mismatched key entry.
func FormatMismatchLine(format OutputFormat, key, left, right string) string {
	switch format {
	case FormatJSON:
		return fmt.Sprintf(`{"type":"mismatch","key":%q,"left":%q,"right":%q}`, key, left, right)
	case FormatMarkdown:
		return fmt.Sprintf("| MISMATCH | `%s` | `%s` | `%s` |"	, key, left, right)
	default:
		return fmt.Sprintf("MISMATCH %-30s left=%-20s right=%s", key, left, right)
	}
}

// MarkdownTableHeader returns the markdown table header string.
func MarkdownTableHeader() string {
	return "| Status | Key | Left | Right |\n|--------|-----|------|-------|\n"
}
