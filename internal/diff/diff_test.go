package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Compare(left, right)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"FOO": "bar", "ONLY_LEFT": "value"}
	right := map[string]string{"FOO": "bar"}

	result := Compare(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no keys missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "bar", "ONLY_RIGHT": "value"}

	result := Compare(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("expected ONLY_RIGHT missing in left, got %v", result.MissingInLeft)
	}
	if len(result.MissingInRight) != 0 {
		t.Errorf("expected no keys missing in right, got %v", result.MissingInRight)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "baz"}

	result := Compare(left, right)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "FOO" || m.LeftValue != "bar" || m.RightValue != "baz" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	left := map[string]string{"ZZZ": "1", "AAA": "1", "MMM": "1"}
	right := map[string]string{}

	result := Compare(left, right)

	expected := []string{"AAA", "MMM", "ZZZ"}
	for i, key := range result.MissingInRight {
		if key != expected[i] {
			t.Errorf("expected sorted key %s at index %d, got %s", expected[i], i, key)
		}
	}
}

func TestCompare_HasDifferences(t *testing.T) {
	empty := Result{}
	if empty.HasDifferences() {
		t.Error("expected HasDifferences to be false for empty result")
	}

	withDiff := Result{MissingInRight: []string{"FOO"}}
	if !withDiff.HasDifferences() {
		t.Error("expected HasDifferences to be true")
	}
}
