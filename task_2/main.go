package main

import (
	"fmt"
	"task_2/wordFrequency"
)

func main() {
	// Example usage
	text1 := "Hello world hello"
	text2 := "lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum!"

	frequency1 := wordFrequency.WordFrequency(text1)
	frequency2 := wordFrequency.WordFrequency(text2)
	fmt.Println(frequency1)
	fmt.Println(frequency2)
}