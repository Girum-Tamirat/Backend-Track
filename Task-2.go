// Fundamentals of Go Tasks
// Task:  Word Frequency Count
// Write a Go function that takes a string as input and returns a dictionary containing the frequency of each word in the string. Treat words in a case-insensitive manner and ignore punctuation marks.
// [Optional]: Write test for your function

// Task : Palindrome Check
// Write a Go function that takes a string as input and checks whether it is a palindrome or not. A palindrome is a word, phrase, number, or other sequence of characters that reads the same forward and backward (ignoring spaces, punctuation, and capitalization).
// [Optional]: Write test for your function


package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// --- Task 1: Word Frequency Count ---
func wordFrequencyCount(text string) map[string]int {
	freq := make(map[string]int)
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove punctuation
	cleaned := ""
	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsSpace(char) {
			cleaned += string(char)
		}
	}

	// Split into words
	words := strings.Fields(cleaned)

	// Count frequency
	for _, word := range words {
		freq[word]++
	}

	return freq
}

// --- Task 2: Palindrome Check ---
func isPalindrome(text string) bool {
	text = strings.ToLower(text)
	cleaned := ""

	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			cleaned += string(char)
		}
	}

	// Compare forward and backward
	n := len(cleaned)
	for i := 0; i < n/2; i++ {
		if cleaned[i] != cleaned[n-1-i] {
			return false
		}
	}
	return true
}

// --- Main Function ---
func main() {
	reader := bufio.NewReader(os.Stdin)

	// Word Frequency Count
	fmt.Print("Enter a phrase for Word Frequency Count: ")
	text, _ := reader.ReadString('\n')
	frequencies := wordFrequencyCount(text)
	fmt.Println("\nWord Frequencies:")
	for word, count := range frequencies {
		fmt.Printf("%s : %d\n", word, count)
	}

	// Palindrome Check
	fmt.Print("\nEnter a text to check if it's a palindrome: ")
	palinText, _ := reader.ReadString('\n')
	if isPalindrome(palinText) {
		fmt.Println("✅ It's a palindrome!")
	} else {
		fmt.Println("❌ Not a palindrome.")
	}
}
