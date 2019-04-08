package plm

import "testing"

func TestIdentity(t *testing.T) {
	reference := Identity{0xab, 0xcd, 0xef}
	value, err := ParseIdentity("abcdef")

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if value != reference {
		t.Errorf("expected:\n%s\ngot:\n%s", reference, value)
	}
}
