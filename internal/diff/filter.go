package diff

// FilterOptions controls which diff results are included in the output.
type FilterOptions struct {
	// OnlyMissing restricts output to variables missing in either environment.
	OnlyMissing bool
	// OnlyMismatched restricts output to variables present in both but with different values.
	OnlyMismatched bool
	// KeyPrefix filters results to only variables whose keys start with the given prefix.
	KeyPrefix string
	// ExcludeKeys is a list of exact key names to omit from the output.
	ExcludeKeys []string
}

// ApplyFilter returns a new Result containing only the entries that pass the
// filter criteria defined in FilterOptions.
func ApplyFilter(r Result, f FilterOptions) Result {
	filtered := Result{}

	if !f.OnlyMismatched {
		for _, entry := range r.MissingInRight {
			if matchesFilter(entry.Key, f) {
				filtered.MissingInRight = append(filtered.MissingInRight, entry)
			}
		}
		for _, entry := range r.MissingInLeft {
			if matchesFilter(entry.Key, f) {
				filtered.MissingInLeft = append(filtered.MissingInLeft, entry)
			}
		}
	}

	if !f.OnlyMissing {
		for _, entry := range r.Mismatched {
			if matchesFilter(entry.Key, f) {
				filtered.Mismatched = append(filtered.Mismatched, entry)
			}
		}
	}

	return filtered
}

// matchesFilter returns true if the key satisfies the prefix and exclusion
// constraints in FilterOptions.
func matchesFilter(key string, f FilterOptions) bool {
	if f.KeyPrefix != "" && !hasPrefix(key, f.KeyPrefix) {
		return false
	}
	if contains(f.ExcludeKeys, key) {
		return false
	}
	return true
}

// hasPrefix reports whether s starts with prefix (case-sensitive).
func hasPrefix(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return s[:len(prefix)] == prefix
}
