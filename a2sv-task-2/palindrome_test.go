package main

import (
	"strings"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"racecar", true},
		{"hello", false},
		{"A man a plan a canal Panama", true},
		{"A man, a plan, a canal--panama", true},
		{"", true},
	}

	for _, test := range tests {
		result := IsPalindrome(NormalizeString(strings.Join(strings.Fields(test.input), "")))
		if result != test.expected {
			t.Errorf("For input '%s', expected %v but got %v", test.input, test.expected, result)
		}
	}
}
