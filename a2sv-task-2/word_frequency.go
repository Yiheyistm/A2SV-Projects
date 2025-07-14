package main

import (
	"strings"
	"unicode"
)

func NormalizeString(s string) string {
    var builder strings.Builder
    for _, ch := range strings.ToLower(s) {
        if !unicode.IsPunct(ch) {
            builder.WriteRune(ch)
        }
    }
    return strings.TrimSpace(builder.String())
}

func CountWordsFrequency(s string) map[string]int {
	var frequency = make(map[string]int)
	words := strings.FieldsSeq(s)
	for word := range words {
		frequency[word]++
	}
	return frequency
}
