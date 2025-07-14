package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func normalizeString(s string) string {
    var builder strings.Builder
    for _, ch := range strings.ToLower(s) {
        if !unicode.IsPunct(ch) {
            builder.WriteRune(ch)
        }
    }
    return strings.TrimSpace(builder.String())
}

func countWordsFrequency(s string) map[string]int {
	var frequency = make(map[string]int)
	words := strings.FieldsSeq(s)
	for word := range words {
		frequency[word]++
	}
	return frequency
}

func main() {
    reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a string: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text) 
    array := strings.Split(text, " ")
	fmt.Println("Array:", array)
	normalizedText := normalizeString(text)
	fmt.Println("Word Frequency:", countWordsFrequency(normalizedText))
   
}