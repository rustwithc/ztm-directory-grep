package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// walkDir recursively lists files and directories
func walkDir(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", path, err)
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			fmt.Printf("[DIR]  %s\n", fullPath)
			walkDir(fullPath) // Recursive call
		} else {
			fmt.Printf("[FILE] %s\n", fullPath)
		}
	}
}

func main() {
		if len(os.Args) > 2 {
		fmt.Println("Supply only one argument: the directory path")
		os.Exit(1)
	}
	startPath := os.Args[1]  // Change to your target directory
	walkDir(startPath)
}
