package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func baseExportResult() *Result {
	return &Result{
		MissingInLeft:  []string{"DB_HOST"},
		MissingInRight: []string{"APP_PORT"},
		Mismatched: []Mismatch{
			{Key: "SECRET_KEY", LeftValue: "abc123", RightValue: "xyz789"},
			{Key: "LOG_LEVEL", LeftValue: "debug", RightValue: "info"},
		},
	}
}

func TestExportAsJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	err := ExportAsJSON(&buf, baseExportResult(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out ExportResult
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(out.MissingInLeft) != 1 || out.MissingInLeft[0] != "DB_HOST" {
		t.Errorf("expected MissingInLeft=[DB_HOST], got %v", out.MissingInLeft)
	}
	if len(out.Mismatched) != 2 {
		t.Errorf("expected 2 mismatched entries, got %d", len(out.Mismatched))
	}
}

func TestExportAsJSON_MasksSecrets(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.MaskSecrets = true
	err := ExportAsJSON(&buf, baseExportResult(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "abc123") {
		t.Error("expected SECRET_KEY value to be masked in JSON output")
	}
	if !strings.Contains(buf.String(), "debug") {
		t.Error("expected non-sensitive LOG_LEVEL value to remain unmasked")
	}
}

func TestExportAsCSV_Headers(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	err := ExportAsCSV(&buf, baseExportResult(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "type,key,left_value,right_value" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if len(lines) != 5 { // header + 1 missing_left + 1 missing_right + 2 mismatched
		t.Errorf("expected 5 CSV lines, got %d", len(lines))
	}
}

func TestExportAsCSV_MissingEntries(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	err := ExportAsCSV(&buf, baseExportResult(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "missing_in_left,DB_HOST") {
		t.Error("expected missing_in_left entry for DB_HOST")
	}
	if !strings.Contains(output, "missing_in_right,APP_PORT") {
		t.Error("expected missing_in_right entry for APP_PORT")
	}
}
