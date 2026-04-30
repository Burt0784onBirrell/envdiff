package diff

import "fmt"

// ValidationIssue represents a single validation problem found in a diff result.
type ValidationIssue struct {
	Key     string
	Message string
	Severity string // "error" or "warning"
}

// ValidationResult holds all issues found during validation.
type ValidationResult struct {
	Issues []ValidationIssue
	Valid  bool
}

// Validate checks a CompareResult against a set of ValidationRules and returns
// a ValidationResult describing any problems found.
func Validate(result CompareResult, rules ValidationRules) ValidationResult {
	var issues []ValidationIssue

	if rules.RequireAllPresent {
		for _, key := range result.MissingInRight {
			issues = append(issues, ValidationIssue{
				Key:      key,
				Message:  fmt.Sprintf("key %q is missing in right file", key),
				Severity: "error",
			})
		}
		for _, key := range result.MissingInLeft {
			issues = append(issues, ValidationIssue{
				Key:      key,
				Message:  fmt.Sprintf("key %q is missing in left file", key),
				Severity: "error",
			})
		}
	}

	if rules.DisallowMismatches {
		for _, mm := range result.Mismatched {
			issues = append(issues, ValidationIssue{
				Key:      mm.Key,
				Message:  fmt.Sprintf("key %q has mismatched values", mm.Key),
				Severity: "warning",
			})
		}
	}

	for _, required := range rules.RequiredKeys {
		if !keyPresentInResult(required, result) {
			issues = append(issues, ValidationIssue{
				Key:      required,
				Message:  fmt.Sprintf("required key %q not found in either file", required),
				Severity: "error",
			})
		}
	}

	return ValidationResult{
		Issues: issues,
		Valid:  len(issues) == 0,
	}
}

func keyPresentInResult(key string, result CompareResult) bool {
	for _, k := range result.MissingInLeft {
		if k == key {
			return true
		}
	}
	for _, k := range result.MissingInRight {
		if k == key {
			return true
		}
	}
	for _, mm := range result.Mismatched {
		if mm.Key == key {
			return true
		}
	}
	return false
}
