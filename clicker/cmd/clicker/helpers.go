package main

import "fmt"

// printCheck prints an actionability check result with a checkmark or X.
func printCheck(name string, passed bool) {
	if passed {
		fmt.Printf("✓ %s: true\n", name)
	} else {
		fmt.Printf("✗ %s: false\n", name)
	}
}
