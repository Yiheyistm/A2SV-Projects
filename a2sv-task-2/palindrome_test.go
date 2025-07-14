package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, test.expected, result, "For input '%s', expected %v but got %v", test.input, test.expected, result)
	}
	
}
