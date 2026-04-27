package diff

import (
	"strings"
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"text", FormatText, false},
		{"TEXT", FormatText, false},
		{"", FormatText, false},
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"markdown", FormatMarkdown, false},
		{"md", FormatMarkdown, false},
		{"MD", FormatMarkdown, false},
		{"xml", FormatText, true},
		{"csv", FormatText, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatMissingLine(t *testing.T) {
	key := "DB_HOST"
	side := "right"

	textOut := FormatMissingLine(FormatText, key, side)
	if !strings.Contains(textOut, "MISSING") || !strings.Contains(textOut, key) {
		t.Errorf("text format missing expected content: %s", textOut)
	}

	jsonOut := FormatMissingLine(FormatJSON, key, side)
	if !strings.Contains(jsonOut, `"type":"missing"`) || !strings.Contains(jsonOut, `"key":"DB_HOST"`) {
		t.Errorf("json format missing expected content: %s", jsonOut)
	}

	mdOut := FormatMissingLine(FormatMarkdown, key, side)
	if !strings.Contains(mdOut, "| MISSING |") || !strings.Contains(mdOut, "`DB_HOST`") {
		t.Errorf("markdown format missing expected content: %s", mdOut)
	}
}

func TestFormatMismatchLine(t *testing.T) {
	key, left, right := "APP_ENV", "development", "production"

	textOut := FormatMismatchLine(FormatText, key, left, right)
	if !strings.Contains(textOut, "MISMATCH") || !strings.Contains(textOut, key) {
		t.Errorf("text format missing expected content: %s", textOut)
	}

	jsonOut := FormatMismatchLine(FormatJSON, key, left, right)
	if !strings.Contains(jsonOut, `"type":"mismatch"`) || !strings.Contains(jsonOut, `"left":"development"`) {
		t.Errorf("json format missing expected content: %s", jsonOut)
	}

	mdOut := FormatMismatchLine(FormatMarkdown, key, left, right)
	if !strings.Contains(mdOut, "| MISMATCH |") || !strings.Contains(mdOut, "`APP_ENV`") {
		t.Errorf("markdown format missing expected content: %s", mdOut)
	}
}

func TestMarkdownTableHeader(t *testing.T) {
	header := MarkdownTableHeader()
	if !strings.Contains(header, "| Status |") {
		t.Errorf("markdown header missing Status column: %s", header)
	}
	if !strings.Contains(header, "|-----") {
		t.Errorf("markdown header missing separator row: %s", header)
	}
}
