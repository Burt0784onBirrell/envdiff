package diff

import (
	"fmt"
	"io"
	"strings"
)

// WriteAnnotationReport writes a human-readable annotation report to the given writer.
func WriteAnnotationReport(w io.Writer, annotations []Annotation) {
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No annotations.")
		return
	}

	counts := map[string]int{}
	for _, a := range annotations {
		counts[a.Severity]++
	}

	fmt.Fprintf(w, "Annotations: %d total", len(annotations))
	parts := []string{}
	for _, sev := range []string{"error", "warning", "info"} {
		if n, ok := counts[sev]; ok {
			parts = append(parts, fmt.Sprintf("%d %s", n, sev))
		}
	}
	if len(parts) > 0 {
		fmt.Fprintf(w, " (%s)", strings.Join(parts, ", "))
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, a := range annotations {
		prefix := annotationPrefix(a.Severity)
		fmt.Fprintf(w, "%s [%s] %s: %s\n", prefix, strings.ToUpper(a.Severity), a.Key, a.Message)
	}
}

func annotationPrefix(severity string) string {
	switch severity {
	case "error":
		return "✖"
	case "warning":
		return "⚠"
	default:
		return "ℹ"
	}
}

// AnnotationSummaryLine returns a one-line summary of the annotation results.
func AnnotationSummaryLine(annotations []Annotation) string {
	if len(annotations) == 0 {
		return "annotations: none"
	}
	counts := map[string]int{}
	for _, a := range annotations {
		counts[a.Severity]++
	}
	parts := []string{}
	for _, sev := range []string{"error", "warning", "info"} {
		if n, ok := counts[sev]; ok {
			parts = append(parts, fmt.Sprintf("%d %s", n, sev))
		}
	}
	return fmt.Sprintf("annotations: %s", strings.Join(parts, ", "))
}
