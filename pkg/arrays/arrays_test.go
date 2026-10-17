package arrays

import "testing"

func TestIn(t *testing.T) {
	values := []string{"one", "two"}
	if !In("two", values) {
		t.Fatal("In() returned false for an existing value")
	}
	if In("three", values) {
		t.Fatal("In() returned true for a missing value")
	}
}
