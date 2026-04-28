package diff

import (
	"encoding/json"
	"fmt"
	"io"
)

// ExportFormat defines supported export formats for diff results.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)

// ExportResult holds the structured diff result for export.
type ExportResult struct {
	MissingInLeft  []string    `json:"missing_in_left,omitempty"`
	MissingInRight []string    `json:"missing_in_right,omitempty"`
	Mismatched     []MismatchEntry `json:"mismatched,omitempty"`
}

// MismatchEntry represents a single mismatched key with both values.
type MismatchEntry struct {
	Key        string `json:"key"`
	LeftValue  string `json:"left_value"`
	RightValue string `json:"right_value"`
}

// ExportJSON writes the diff result as JSON to the given writer.
func ExportAsJSON(w io.Writer, result *Result, opts ReportOptions) error {
	export := buildExportResult(result, opts)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(export)
}

// ExportAsCSV writes the diff result as CSV to the given writer.
func ExportAsCSV(w io.Writer, result *Result, opts ReportOptions) error {
	fmt.Fprintln(w, "type,key,left_value,right_value")
	for _, key := range result.MissingInLeft {
		fmt.Fprintf(w, "missing_in_left,%s,,\n", key)
	}
	for _, key := range result.MissingInRight {
		fmt.Fprintf(w, "missing_in_right,%s,,\n", key)
	}
	for _, mm := range result.Mismatched {
		lv := maskIfSensitive(mm.Key, mm.LeftValue, opts)
		rv := maskIfSensitive(mm.Key, mm.RightValue, opts)
		fmt.Fprintf(w, "mismatched,%s,%s,%s\n", mm.Key, lv, rv)
	}
	return nil
}

func buildExportResult(result *Result, opts ReportOptions) ExportResult {
	entries := make([]MismatchEntry, 0, len(result.Mismatched))
	for _, mm := range result.Mismatched {
		entries = append(entries, MismatchEntry{
			Key:        mm.Key,
			LeftValue:  maskIfSensitive(mm.Key, mm.LeftValue, opts),
			RightValue: maskIfSensitive(mm.Key, mm.RightValue, opts),
		})
	}
	return ExportResult{
		MissingInLeft:  result.MissingInLeft,
		MissingInRight: result.MissingInRight,
		Mismatched:     entries,
	}
}

func maskIfSensitive(key, value string, opts ReportOptions) string {
	if opts.MaskSecrets && IsSensitiveKey(key, opts.SensitiveKeys) {
		return maskValue(value)
	}
	return value
}
