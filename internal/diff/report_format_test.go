package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteFormattedReport_Plain_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.Format = FormatPlain
	err := WriteFormattedReport(&buf, &Result{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWriteFormattedReport_Plain_WithDiff(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.Format = FormatPlain
	result := &Result{
		MissingInRight: []string{"API_KEY"},
		Mismatched:     []Mismatch{{Key: "LOG_LEVEL", LeftValue: "debug", RightValue: "warn"}},
	}
	err := WriteFormattedReport(&buf, result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "API_KEY") {
		t.Error("expected API_KEY in plain output")
	}
	if !strings.Contains(output, "Summary:") {
		t.Error("expected Summary line in plain output")
	}
}

func TestWriteFormattedReport_Markdown_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.Format = FormatMarkdown
	err := WriteFormattedReport(&buf, &Result{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "_No differences found._") {
		t.Errorf("expected markdown no-diff message, got: %s", buf.String())
	}
}

func TestWriteFormattedReport_Markdown_WithDiff(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.Format = FormatMarkdown
	result := &Result{
		MissingInLeft: []string{"DB_PASS"},
	}
	err := WriteFormattedReport(&buf, result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "|") {
		t.Error("expected markdown table in output")
	}
	if !strings.Contains(output, "> Summary:") {
		t.Error("expected markdown summary line")
	}
}

func TestWriteFormattedReport_JSON(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.Format = FormatJSON
	result := &Result{MissingInRight: []string{"REDIS_URL"}}
	err := WriteFormattedReport(&buf, result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "REDIS_URL") {
		t.Error("expected REDIS_URL in JSON output")
	}
}

func TestWriteFormattedReport_CSV(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions([]string{}, []string{})
	opts.Format = FormatCSV
	result := &Result{MissingInLeft: []string{"TIMEOUT"}}
	err := WriteFormattedReport(&buf, result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "TIMEOUT") {
		t.Error("expected TIMEOUT in CSV output")
	}
}
