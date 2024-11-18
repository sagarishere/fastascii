package main

import (
	"fmt"
	"os"
	"strings"
)

var myMap map[rune][]string

func main() {
	mybytes, _ := os.ReadFile("standard.txt")
	myMap = fillMap(squeeze(strings.Split(strings.ReplaceAll(string(mybytes), "\\r\\n", "\\n"), "\n")))
	myArgs := getArgs()
	words := strings.Split(myArgs, "\\n")
	printWords(words)
}

func printWords(words []string) {
	for _, word := range words {
		if word == "" {
			fmt.Println()
		} else {
			printWord(word)
		}
	}
}

func printWord(word string) {
	finalPrint := ""
	for i := 0; i < 8; i++ {
		for _, rune := range word {
			finalPrint += myMap[rune][i]
		}
		finalPrint += "\n"
	}
	fmt.Print(finalPrint)
}

func getArgs() string {
	myArgs := os.Args[1:]
	if len(myArgs) != 1 {
		fmt.Println("Usage: go run main.go <words>")
		os.Exit(1)
	}
	return myArgs[0]
}

func fillMap(squeezed []string) map[rune][]string {
	result := make(map[rune][]string)
	for i := 0; i < len(squeezed); i += 8 {
		result[(rune(i/8 + ' '))] = squeezed[i : i+8]
	}
	return result
}

func squeeze(content []string) (result []string) {
	for _, line := range content {
		if line == "" {
			continue
		}
		result = append(result, line)
	}
	return
}
