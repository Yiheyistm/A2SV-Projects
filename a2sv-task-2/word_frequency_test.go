package main

import "testing"

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "hello world"},
		{"  GoLang  ", "golang"},
		{"Test@123", "test123"},
	}

	for _, test := range tests {
		result := NormalizeString(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s' but got '%s'", test.input, test.expected, result)
		}
	}
}

func TestCountWordsFrequency(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]int
	}{
		{"hello world hello", map[string]int{"hello": 2, "world": 1}},
		{"Go Go Go", map[string]int{"go": 3}},
		{"", map[string]int{}},
	}

	for _, test := range tests {
		result := CountWordsFrequency(test.input)
		for word, count := range test.expected {
			if result[word] != count {
				t.Errorf("For word '%s' in input '%s', expected %d but got %d", word, test.input, count, result[word])
			}
		}
	}
}
