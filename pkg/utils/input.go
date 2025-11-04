package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptString prompts the user for input and returns the trimmed string
func PromptString(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// PromptYesNo prompts the user for yes/no input
func PromptYesNo(prompt string, defaultYes bool) (bool, error) {
	defaultStr := "n"
	if defaultYes {
		defaultStr = "y"
	}

	fmt.Printf("%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultYes, nil
	}

	return input == "y" || input == "yes", nil
}

// PromptChoice prompts the user to select from a list of choices
func PromptChoice(prompt string, choices []string) (int, error) {
	fmt.Println(prompt)
	for i, choice := range choices {
		fmt.Printf("  %d. %s\n", i+1, choice)
	}

	fmt.Print("Select (number): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return -1, fmt.Errorf("no selection made")
	}

	var selection int
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil {
		return -1, fmt.Errorf("invalid selection: %s", input)
	}

	if selection < 1 || selection > len(choices) {
		return -1, fmt.Errorf("selection out of range")
	}

	return selection - 1, nil
}
