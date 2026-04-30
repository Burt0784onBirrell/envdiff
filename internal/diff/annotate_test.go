package diff

import (
	"bytes"
	"strings"
	"testing"
)

func baseAnnotateResult() CompareResult {
	return CompareResult{
		MissingInRight: []string{"DB_PASSWORD", "APP_NAME"},
		MissingInLeft:  []string{"SECRET_KEY"},
		Mismatched: []MismatchedVar{
			{Key: "DATABASE_URL", LeftValue: "postgres://localhost", RightValue: "postgres://prod"},
		},
	}
}

func TestAnnotate_SensitiveKeys(t *testing.T) {
	result := baseAnnotateResult()
	annotations := Annotate(result, DefaultAnnotationRules())

	keys := map[string]bool{}
	for _, a := range annotations {
		keys[a.Key] = true
	}

	if !keys["DB_PASSWORD"] {
		t.Error("expected annotation for DB_PASSWORD")
	}
	if !keys["SECRET_KEY"] {
		t.Error("expected annotation for SECRET_KEY")
	}
}

func TestAnnotate_DatabaseKeys(t *testing.T) {
	result := baseAnnotateResult()
	annotations := Annotate(result, DefaultAnnotationRules())

	var dbAnn *Annotation
	for _, a := range annotations {
		if a.Key == "DATABASE_URL" && a.Severity == "error" {
			copy := a
			dbAnn = &copy
			break
		}
	}
	if dbAnn == nil {
		t.Error("expected error annotation for DATABASE_URL")
	}
}

func TestAnnotate_NoIssues(t *testing.T) {
	result := CompareResult{}
	annotations := Annotate(result, DefaultAnnotationRules())
	if len(annotations) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(annotations))
	}
}

func TestWriteAnnotationReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteAnnotationReport(&buf, []Annotation{})
	if !strings.Contains(buf.String(), "No annotations") {
		t.Error("expected 'No annotations' message")
	}
}

func TestWriteAnnotationReport_WithAnnotations(t *testing.T) {
	annotations := []Annotation{
		{Key: "SECRET_KEY", Message: "sensitive key", Severity: "warning"},
		{Key: "DATABASE_URL", Message: "db key", Severity: "error"},
	}
	var buf bytes.Buffer
	WriteAnnotationReport(&buf, annotations)
	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Error("expected SECRET_KEY in report")
	}
	if !strings.Contains(out, "DATABASE_URL") {
		t.Error("expected DATABASE_URL in report")
	}
	if !strings.Contains(out, "WARNING") {
		t.Error("expected WARNING severity label")
	}
}

func TestAnnotationSummaryLine(t *testing.T) {
	annotations := []Annotation{
		{Key: "K1", Severity: "error"},
		{Key: "K2", Severity: "warning"},
		{Key: "K3", Severity: "warning"},
	}
	line := AnnotationSummaryLine(annotations)
	if !strings.Contains(line, "1 error") {
		t.Errorf("expected '1 error' in summary, got: %s", line)
	}
	if !strings.Contains(line, "2 warning") {
		t.Errorf("expected '2 warning' in summary, got: %s", line)
	}
}
