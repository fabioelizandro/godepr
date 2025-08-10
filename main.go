package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		printHelp()
		return
	}

	// Read go list JSON from stdin
	jsonData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON data
	packages, err := ParsePackagesJSON(string(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing packages: %v\n", err)
		os.Exit(1)
	}

	// For now, create a sample rule (later this will come from .godepr file)
	// This is just for demonstration
	rules := []Rule{
		{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		},
	}

	// Check rules
	violations := CheckRules(packages, rules)

	// Report violations
	if len(violations) > 0 {
		fmt.Println("Dependency violations found:")
		fmt.Println()
		
		for _, v := range violations {
			fmt.Printf("  %s\n", v.Message)
		}
		
		fmt.Println()
		fmt.Printf("Total violations: %d\n", len(violations))
		os.Exit(1)
	}

	fmt.Println("No dependency violations found")
}

func printHelp() {
	fmt.Println(`godepr - Go Dependency Governance Tool

Usage:
  go list -json ./... | godepr

Description:
  godepr reads the output of 'go list -json' from stdin and checks
  if the imports comply with the rules defined in .godepr file.

Example:
  # Check all packages in current directory
  go list -json ./... | godepr

  # Check specific directory
  go list -json ./internal/... | godepr

Future Features:
  - Read rules from .godepr configuration file
  - Support for multiple rule types (allowed-list, std-dep-only, etc.)
  - Better error reporting with suggestions

Current Implementation:
  - Supports denied-list rule type
  - Reads go list JSON from stdin
  - Reports violations to stdout`)
}