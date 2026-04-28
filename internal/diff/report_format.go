package diff

import (
	"fmt"
	"io"
)

// WriteFormattedReport writes a diff report using the specified output format.
func WriteFormattedReport(w io.Writer, result *Result, opts ReportOptions) error {
	switch opts.Format {
	case FormatMarkdown:
		return writeMarkdownReport(w, result, opts)
	case FormatJSON:
		return ExportAsJSON(w, result, opts)
	case FormatCSV:
		return ExportAsCSV(w, result, opts)
	default:
		return writePlainReport(w, result, opts)
	}
}

func writePlainReport(w io.Writer, result *Result, opts ReportOptions) error {
	if len(result.MissingInLeft) == 0 && len(result.MissingInRight) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}
	for _, key := range result.MissingInLeft {
		fmt.Fprintln(w, FormatMissingLine(key, "left", opts.Format))
	}
	for _, key := range result.MissingInRight {
		fmt.Fprintln(w, FormatMissingLine(key, "right", opts.Format))
	}
	for _, mm := range result.Mismatched {
		lv := maskIfSensitive(mm.Key, mm.LeftValue, opts)
		rv := maskIfSensitive(mm.Key, mm.RightValue, opts)
		fmt.Fprintln(w, FormatMismatchLine(mm.Key, lv, rv, opts.Format))
	}
	fmt.Fprintln(w, SummaryLine(result))
	return nil
}

func writeMarkdownReport(w io.Writer, result *Result, opts ReportOptions) error {
	if len(result.MissingInLeft) == 0 && len(result.MissingInRight) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, "_No differences found._")
		return nil
	}
	fmt.Fprintln(w, MarkdownTableHeader())
	for _, key := range result.MissingInLeft {
		fmt.Fprintln(w, FormatMissingLine(key, "left", opts.Format))
	}
	for _, key := range result.MissingInRight {
		fmt.Fprintln(w, FormatMissingLine(key, "right", opts.Format))
	}
	for _, mm := range result.Mismatched {
		lv := maskIfSensitive(mm.Key, mm.LeftValue, opts)
		rv := maskIfSensitive(mm.Key, mm.RightValue, opts)
		fmt.Fprintln(w, FormatMismatchLine(mm.Key, lv, rv, opts.Format))
	}
	fmt.Fprintf(w, "\n> %s\n", SummaryLine(result))
	return nil
}

// SummaryLine returns a human-readable summary of the diff result.
func SummaryLine(result *Result) string {
	return fmt.Sprintf(
		"Summary: %d missing in left, %d missing in right, %d mismatched.",
		len(result.MissingInLeft),
		len(result.MissingInRight),
		len(result.Mismatched),
	)
}
