package unixid_test

import (
	"testing"

	"github.com/cdvelop/unixid"
)

func TestValidateID(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
		err      bool
	}{
		{"1624397134562544800", 1624397134562544800, false},
		{"1624397134562544800.42", 1624397134562544800, false},
		{"1624397134562544800.42.42", 0, true},
		{"1624397134562544800a", 0, true},
		{".1624397134562544800", 0, true},
		{"1624397134562544800.", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			id, err := unixid.ValidateID(tc.input)
			if (err != nil) != tc.err {
				t.Errorf("expected error: %v, got: %v", tc.err, err)
			}
			if id != tc.expected {
				t.Errorf("expected: %d, got: %d", tc.expected, id)
			}
		})
	}
}
