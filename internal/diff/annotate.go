package diff

import "fmt"

// Annotation represents a human-readable note attached to a diff key.
type Annotation struct {
	Key     string
	Message string
	Severity string // "info", "warning", "error"
}

// AnnotationRule defines a function that may produce an annotation for a given key and value.
type AnnotationRule func(key, value string) *Annotation

// Annotate applies a set of annotation rules to the keys in a CompareResult
// and returns a slice of Annotations for any matched keys.
func Annotate(result CompareResult, rules []AnnotationRule) []Annotation {
	var annotations []Annotation

	allKeys := make(map[string]string)
	for _, k := range result.MissingInRight {
		allKeys[k] = ""
	}
	for _, k := range result.MissingInLeft {
		allKeys[k] = ""
	}
	for _, m := range result.Mismatched {
		allKeys[m.Key] = m.LeftValue
	}

	for key, val := range allKeys {
		for _, rule := range rules {
			if ann := rule(key, val); ann != nil {
				annotations = append(annotations, *ann)
			}
		}
	}

	return annotations
}

// DefaultAnnotationRules returns a standard set of annotation rules.
func DefaultAnnotationRules() []AnnotationRule {
	return []AnnotationRule{
		func(key, _ string) *Annotation {
			if IsSensitiveKey(key) {
				return &Annotation{
					Key:      key,
					Message:  fmt.Sprintf("%s appears to be a sensitive key — verify it is set correctly in all environments", key),
					Severity: "warning",
				}
			}
			return nil
		},
		func(key, _ string) *Annotation {
			if hasPrefix(toUpper(key), "DATABASE") || hasPrefix(toUpper(key), "DB_") {
				return &Annotation{
					Key:      key,
					Message:  fmt.Sprintf("%s is a database-related key — mismatches may cause connectivity issues", key),
					Severity: "error",
				}
			}
			return nil
		},
	}
}
