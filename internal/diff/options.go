package diff

// ReportOptions controls how the diff report is rendered.
type ReportOptions struct {
	// LeftLabel is the display name for the left (base) env file.
	LeftLabel string
	// RightLabel is the display name for the right (target) env file.
	RightLabel string
	// MaskSecrets redacts values for keys that look sensitive.
	MaskSecrets bool
}

// DefaultReportOptions returns sensible defaults for report generation.
func DefaultReportOptions(leftPath, rightPath string) ReportOptions {
	label := func(p string) string {
		if p == "" {
			return "<unknown>"
		}
		return p
	}
	return ReportOptions{
		LeftLabel:   label(leftPath),
		RightLabel:  label(rightPath),
		MaskSecrets: true,
	}
}

// sensitiveKeyPatterns holds substrings that indicate a key holds a secret.
var sensitiveKeyPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// IsSensitiveKey returns true when the key name suggests it holds a secret
// value that should be masked in reports.
func IsSensitiveKey(key string) bool {
	upper := toUpper(key)
	for _, pattern := range sensitiveKeyPatterns {
		if contains(upper, pattern) {
			return true
		}
	}
	return false
}

func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
