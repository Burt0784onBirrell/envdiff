package diff

import (
	"fmt"
	"io"
)

// WriteValidationReport writes a human-readable summary of validation issues to w.
func WriteValidationReport(w io.Writer, vr ValidationResult) {
	if vr.Valid {
		fmt.Fprintln(w, "✔ Validation passed: no issues found.")
		return
	}

	errors := 0
	warnings := 0
	for _, issue := range vr.Issues {
		switch issue.Severity {
		case "error":
			errors++
		case "warning":
			warnings++
		}
	}

	fmt.Fprintf(w, "✘ Validation failed: %d error(s), %d warning(s)\n\n", errors, warnings)

	if errors > 0 {
		fmt.Fprintln(w, "Errors:")
		for _, issue := range vr.Issues {
			if issue.Severity == "error" {
				fmt.Fprintf(w, "  [ERROR] %s: %s\n", issue.Key, issue.Message)
			}
		}
	}

	if warnings > 0 {
		fmt.Fprintln(w, "Warnings:")
		for _, issue := range vr.Issues {
			if issue.Severity == "warning" {
				fmt.Fprintf(w, "  [WARN]  %s: %s\n", issue.Key, issue.Message)
			}
		}
	}
}

// ValidationSummaryLine returns a one-line summary suitable for CLI output.
func ValidationSummaryLine(vr ValidationResult) string {
	if vr.Valid {
		return "validation: OK"
	}
	return fmt.Sprintf("validation: FAILED (%d errors, %d warnings)",
		vr.ErrorCount(), countBySeverity(vr, "warning"))
}

func countBySeverity(vr ValidationResult, severity string) int {
	count := 0
	for _, issue := range vr.Issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}
