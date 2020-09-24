package pdf

import (
	"fmt"
	"testing"
)

func TestNewBasicType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		input    string
	}{
		{
			name:     "boolean basic type",
			expected: "true",
			input:    "true",
		},
		{
			name:     "string basic type",
			expected: "hello(1), world",
			input:    "(hello\\(1\\), \\\nworld)",
		},
		{
			name:     "integer basic type",
			expected: "324",
			input:    "+324",
		},
		{
			name:     "real basic type",
			expected: "-0.23",
			input:    "-.23",
		},
		{
			name:     "hexadecimal basic type",
			expected: "4E6F76",
			input:    "<4E6F76>",
		},
		{
			name:     "name basic type",
			expected: "test",
			input:    "/test",
		},
		{
			name:     "array basic type",
			expected: "[test -0.23 [3 you] true]",
			input:    "[/test -.23 [ 3 (you) ] true]",
		},
	}

	for i, test := range tests {
		val := NewBasicType(test.input)

		if test.expected != fmt.Sprintf("%s", val) {
			t.Errorf("Case %d (%s) Failed:\n\tExpected: %v\n\tGot: %s\n",
				i, test.name, test.expected, val)
		}
	}
}
