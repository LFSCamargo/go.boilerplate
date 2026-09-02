package requests

import "testing"

func TestEmailSchema_RejectsInvalid(t *testing.T) {
	var body Email
	if issues := ParseJSON([]byte(`{"email":"nope"}`), EmailSchema, &body); issues == nil {
		t.Fatal("expected issues")
	}
}

func TestEmailSchema_AcceptsValid(t *testing.T) {
	var body Email
	if issues := ParseJSON([]byte(`{"email":"ada@example.com"}`), EmailSchema, &body); issues != nil {
		t.Fatalf("unexpected issues: %#v", Flatten(issues))
	}
}
