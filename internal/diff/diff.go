package diff

// Result holds the comparison result between two env files.
type Result struct {
	// MissingInRight contains keys present in left but not in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but not in left.
	MissingInLeft []string
	// Mismatched contains keys present in both but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey represents a key whose value differs between two env files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare compares two parsed env maps and returns a Result describing
// missing and mismatched keys.
func Compare(left, right map[string]string) Result {
	result := Result{}

	for key, leftVal := range left {
		rightVal, ok := right[key]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, key)
			continue
		}
		if leftVal != rightVal {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        key,
				LeftValue:  leftVal,
				RightValue: rightVal,
			})
		}
	}

	for key := range right {
		if _, ok := left[key]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, key)
		}
	}

	sortStrings(result.MissingInRight)
	sortStrings(result.MissingInLeft)
	sortMismatched(result.Mismatched)

	return result
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInRight) > 0 || len(r.MissingInLeft) > 0 || len(r.Mismatched) > 0
}
