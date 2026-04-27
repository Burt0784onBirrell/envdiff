package diff

import "testing"

func TestDefaultReportOptions(t *testing.T) {
	opts := DefaultReportOptions(".env.local", ".env.production")
	if opts.LeftLabel != ".env.local" {
		t.Errorf("expected LeftLabel '.env.local', got %q", opts.LeftLabel)
	}
	if opts.RightLabel != ".env.production" {
		t.Errorf("expected RightLabel '.env.production', got %q", opts.RightLabel)
	}
	if !opts.MaskSecrets {
		t.Error("expected MaskSecrets to default to true")
	}
}

func TestDefaultReportOptions_EmptyPaths(t *testing.T) {
	opts := DefaultReportOptions("", "")
	if opts.LeftLabel != "<unknown>" {
		t.Errorf("expected '<unknown>', got %q", opts.LeftLabel)
	}
	if opts.RightLabel != "<unknown>" {
		t.Errorf("expected '<unknown>', got %q", opts.RightLabel)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key       string
		sensitive bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"SECRET_KEY", true},
		{"PRIVATE_KEY", true},
		{"AWS_SECRET_ACCESS_KEY", true},
		{"DB_HOST", false},
		{"APP_ENV", false},
		{"PORT", false},
		{"LOG_LEVEL", false},
		{"STRIPE_TOKEN", true},
		{"USER_CREDENTIAL", true},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := IsSensitiveKey(tc.key)
			if got != tc.sensitive {
				t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.sensitive)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	if got := toUpper("hello_world"); got != "HELLO_WORLD" {
		t.Errorf("toUpper(\"hello_world\") = %q, want %q", got, "HELLO_WORLD")
	}
}

func TestContains(t *testing.T) {
	if !contains("HELLO_WORLD", "WORLD") {
		t.Error("expected contains to return true")
	}
	if contains("HELLO", "WORLD") {
		t.Error("expected contains to return false")
	}
}
