package responses

import "testing"

func TestHealthSchema_AcceptsOK(t *testing.T) {
	resp := OK()
	if issues := HealthSchema.Validate(&resp); issues != nil {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}

func TestHealthSchema_RejectsOtherStatus(t *testing.T) {
	resp := Health{Status: "degraded"}
	if issues := HealthSchema.Validate(&resp); issues == nil {
		t.Fatal("expected validation issues")
	}
}
