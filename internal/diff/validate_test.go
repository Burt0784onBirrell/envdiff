package diff

import (
	"testing"
)

func baseValidateResult() CompareResult {
	return CompareResult{
		MissingInLeft:  []string{"DB_HOST"},
		MissingInRight: []string{"API_KEY"},
		Mismatched: []MismatchedVar{
			{Key: "PORT", LeftValue: "8080", RightValue: "9090"},
		},
	}
}

func TestValidate_AllPresent_Fails(t *testing.T) {
	rules := DefaultValidationRules()
	vr := Validate(baseValidateResult(), rules)

	if vr.Valid {
		t.Error("expected validation to fail")
	}
	if !vr.HasErrors() {
		t.Error("expected at least one error")
	}
	if vr.ErrorCount() != 2 {
		t.Errorf("expected 2 errors, got %d", vr.ErrorCount())
	}
}

func TestValidate_NoMismatches_Warning(t *testing.T) {
	rules := StrictValidationRules()
	vr := Validate(baseValidateResult(), rules)

	if !vr.HasWarnings() {
		t.Error("expected at least one warning for mismatched values")
	}
}

func TestValidate_Clean_Result(t *testing.T) {
	rules := DefaultValidationRules()
	clean := CompareResult{}
	vr := Validate(clean, rules)

	if !vr.Valid {
		t.Error("expected validation to pass for empty result")
	}
	if len(vr.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(vr.Issues))
	}
}

func TestValidate_RequiredKeys_Missing(t *testing.T) {
	rules := ValidationRules{
		RequireAllPresent:  false,
		DisallowMismatches: false,
		RequiredKeys:       []string{"MUST_EXIST"},
	}
	vr := Validate(baseValidateResult(), rules)

	if vr.Valid {
		t.Error("expected validation to fail due to missing required key")
	}
	if vr.ErrorCount() != 1 {
		t.Errorf("expected 1 error for missing required key, got %d", vr.ErrorCount())
	}
}

func TestValidate_RequiredKeys_Present(t *testing.T) {
	rules := ValidationRules{
		RequireAllPresent:  false,
		DisallowMismatches: false,
		RequiredKeys:       []string{"PORT"},
	}
	vr := Validate(baseValidateResult(), rules)

	if !vr.Valid {
		t.Errorf("expected validation to pass; issues: %+v", vr.Issues)
	}
}

func TestHasErrors_HasWarnings(t *testing.T) {
	vr := ValidationResult{
		Issues: []ValidationIssue{
			{Key: "A", Severity: "error", Message: "missing"},
			{Key: "B", Severity: "warning", Message: "mismatch"},
		},
	}
	if !vr.HasErrors() {
		t.Error("expected HasErrors to return true")
	}
	if !vr.HasWarnings() {
		t.Error("expected HasWarnings to return true")
	}
}
