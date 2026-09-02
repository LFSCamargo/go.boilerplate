package middleware

import "testing"

func TestEmptySchema_AcceptsZeroValue(t *testing.T) {
	var body Empty
	if issues := EmptySchema.Validate(&body); issues != nil {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}
