package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func main() {
	fmt.Println("Task 2: Word Frequency")
    reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a string: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text) 

	normalizedText := NormalizeString(text)
	fmt.Println("Word Frequency:", CountWordsFrequency(normalizedText))

	fmt.Println("\n===========================================")

	fmt.Println("Task 2: Palindrome Checker")
	fmt.Print("Enter a string: ")
	text2, _ := reader.ReadString('\n')
	text2 = strings.TrimSpace(text2)
	text2 = NormalizeString(text2)
	text2 = strings.Join(strings.Fields(text2), "")

	if IsPalindrome(text2) {
		fmt.Println("The string is a palindrome.")
	} else {
		fmt.Println("The string is not a palindrome.")
	}

}