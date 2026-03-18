package assets

import (
	"fmt"
	"os"
)

func PrintLogo() {
	// Define the file path. Use a raw string literal (backticks) for Windows paths
	// to avoid issues with escape sequences.
	filePath := "./assets/logo.txt"

	// Read the entire file content into a byte slice
	data, err := os.ReadFile(filePath)

	// Error handling is crucial in Go
	if err != nil {
		fmt.Printf("File reading error: %v\n", err)
		return
	}

	// Convert the byte slice to a string and print it to standard output
	fmt.Println(string(data))
}
