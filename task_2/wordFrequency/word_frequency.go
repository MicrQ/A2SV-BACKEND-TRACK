package wordFrequency

import (
	"strings"
	"unicode"
)

// counts the frequency of each word in the given text
func WordFrequency(text string) map[string]int {
	frq := make(map[string]int)

	text = removePunctuation(text) // remove punctuations

	// splitting words by spaces and ignoring punctuation
	words := strings.Fields(text)
	for _, word := range words {
		frq[strings.ToLower(word)]++
	}
	return frq
}

// helper to remove punctuation from the text
func removePunctuation(text string) string {
	var result string

	for _, char := range text {
		if !unicode.IsPunct(char) {
			result += string(char)
		}
	}

	return result
}
