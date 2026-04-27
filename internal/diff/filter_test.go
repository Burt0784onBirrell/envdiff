package diff

import (
	"testing"
)

func baseResult() Result {
	return Result{
		MissingInRight: []Entry{{Key: "DB_HOST", Value: "localhost"}, {Key: "APP_PORT", Value: "8080"}},
		MissingInLeft:  []Entry{{Key: "REDIS_URL", Value: "redis://localhost"}},
		Mismatched:     []MismatchEntry{{Key: "SECRET_KEY", LeftValue: "abc", RightValue: "xyz"}, {Key: "APP_ENV", LeftValue: "dev", RightValue: "prod"}},
	}
}

func TestApplyFilter_NoFilter(t *testing.T) {
	r := ApplyFilter(baseResult(), FilterOptions{})
	if len(r.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(r.MissingInRight))
	}
	if len(r.MissingInLeft) != 1 {
		t.Errorf("expected 1 MissingInLeft, got %d", len(r.MissingInLeft))
	}
	if len(r.Mismatched) != 2 {
		t.Errorf("expected 2 Mismatched, got %d", len(r.Mismatched))
	}
}

func TestApplyFilter_OnlyMissing(t *testing.T) {
	r := ApplyFilter(baseResult(), FilterOptions{OnlyMissing: true})
	if len(r.MissingInRight) != 2 {
		t.Errorf("expected 2 MissingInRight, got %d", len(r.MissingInRight))
	}
	if len(r.Mismatched) != 0 {
		t.Errorf("expected 0 Mismatched, got %d", len(r.Mismatched))
	}
}

func TestApplyFilter_OnlyMismatched(t *testing.T) {
	r := ApplyFilter(baseResult(), FilterOptions{OnlyMismatched: true})
	if len(r.Mismatched) != 2 {
		t.Errorf("expected 2 Mismatched, got %d", len(r.Mismatched))
	}
	if len(r.MissingInRight) != 0 {
		t.Errorf("expected 0 MissingInRight, got %d", len(r.MissingInRight))
	}
}

func TestApplyFilter_KeyPrefix(t *testing.T) {
	r := ApplyFilter(baseResult(), FilterOptions{KeyPrefix: "APP_"})
	if len(r.MissingInRight) != 1 || r.MissingInRight[0].Key != "APP_PORT" {
		t.Errorf("expected only APP_PORT in MissingInRight, got %+v", r.MissingInRight)
	}
	if len(r.Mismatched) != 1 || r.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("expected only APP_ENV in Mismatched, got %+v", r.Mismatched)
	}
}

func TestApplyFilter_ExcludeKeys(t *testing.T) {
	r := ApplyFilter(baseResult(), FilterOptions{ExcludeKeys: []string{"SECRET_KEY", "DB_HOST"}})
	if len(r.MissingInRight) != 1 || r.MissingInRight[0].Key != "APP_PORT" {
		t.Errorf("expected only APP_PORT in MissingInRight, got %+v", r.MissingInRight)
	}
	if len(r.Mismatched) != 1 || r.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("expected only APP_ENV in Mismatched, got %+v", r.Mismatched)
	}
}

func TestHasPrefix(t *testing.T) {
	cases := []struct {
		s, prefix string
		want     bool
	}{
		{"APP_PORT", "APP_", true},
		{"DB_HOST", "APP_", false},
		{"APP", "APP_", false},
		{"APP_", "APP_", true},
	}
	for _, c := range cases {
		if got := hasPrefix(c.s, c.prefix); got != c.want {
			t.Errorf("hasPrefix(%q, %q) = %v, want %v", c.s, c.prefix, got, c.want)
		}
	}
}
