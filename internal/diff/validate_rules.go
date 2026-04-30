package diff

// ValidationRules defines the constraints applied during validation of a diff result.
type ValidationRules struct {
	// RequireAllPresent fails if any key is missing from either side.
	RequireAllPresent bool

	// DisallowMismatches warns when values differ between the two files.
	DisallowMismatches bool

	// RequiredKeys is a list of keys that must appear in at least one of the files.
	RequiredKeys []string
}

// DefaultValidationRules returns a ValidationRules with sensible defaults.
func DefaultValidationRules() ValidationRules {
	return ValidationRules{
		RequireAllPresent:  true,
		DisallowMismatches: false,
		RequiredKeys:       []string{},
	}
}

// StrictValidationRules returns a ValidationRules that treats all differences as errors.
func StrictValidationRules() ValidationRules {
	return ValidationRules{
		RequireAllPresent:  true,
		DisallowMismatches: true,
		RequiredKeys:       []string{},
	}
}

// HasErrors reports whether the ValidationResult contains any error-severity issues.
func (vr ValidationResult) HasErrors() bool {
	for _, issue := range vr.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// HasWarnings reports whether the ValidationResult contains any warning-severity issues.
func (vr ValidationResult) HasWarnings() bool {
	for _, issue := range vr.Issues {
		if issue.Severity == "warning" {
			return true
		}
	}
	return false
}

// ErrorCount returns the number of error-severity issues.
func (vr ValidationResult) ErrorCount() int {
	count := 0
	for _, issue := range vr.Issues {
		if issue.Severity == "error" {
			count++
		}
	}
	return count
}
