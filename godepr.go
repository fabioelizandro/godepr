package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Rule struct {
	Directory string   `json:"directory"`
	RuleType  string   `json:"ruletype"`
	RuleBody  []string `json:"rulebody"`
}

type Config struct {
	Rules []Rule `json:"rules"`
}

type Package struct {
	ImportPath string   `json:"ImportPath"`
	Dir        string   `json:"Dir"`
	Deps       []string `json:"Deps"`
	Imports    []string `json:"Imports"`
}

type Violation struct {
	PackagePath string
	ImportPath  string
	RuleType    string
	Message     string
}

func CheckDeniedListRule(packages []Package, rule Rule) []Violation {
	var violations []Violation
	
	targetDir := normalizeDirectory(rule.Directory)
	
	for _, pkg := range packages {
		if !isInDirectory(pkg.Dir, targetDir) {
			continue
		}
		
		for _, imp := range pkg.Imports {
			for _, pattern := range rule.RuleBody {
				if matchesPattern(imp, pattern) {
					violations = append(violations, Violation{
						PackagePath: pkg.ImportPath,
						ImportPath:  imp,
						RuleType:    rule.RuleType,
						Message:     fmt.Sprintf("package %s imports denied dependency %s", pkg.ImportPath, imp),
					})
					break
				}
			}
		}
	}
	
	return violations
}

func CheckRules(packages []Package, rules []Rule) []Violation {
	var allViolations []Violation
	
	for _, rule := range rules {
		switch rule.RuleType {
		case "denied-list":
			violations := CheckDeniedListRule(packages, rule)
			allViolations = append(allViolations, violations...)
		default:
			// Unsupported rule type, skip
		}
	}
	
	return allViolations
}

func ParsePackagesJSON(jsonData string) ([]Package, error) {
	var packages []Package
	decoder := json.NewDecoder(strings.NewReader(jsonData))
	
	for decoder.More() {
		var pkg Package
		if err := decoder.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, pkg)
	}
	
	return packages, nil
}

func normalizeDirectory(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return dir
}

func isInDirectory(pkgDir, targetDir string) bool {
	// Normalize paths for comparison
	if !strings.HasSuffix(pkgDir, "/") {
		pkgDir += "/"
	}
	
	// Check if pkgDir is within targetDir
	// This handles both exact matches and subdirectories
	return strings.HasPrefix(pkgDir, targetDir) || 
		strings.HasSuffix(pkgDir, "/"+strings.TrimSuffix(targetDir, "/")+"/") ||
		strings.Contains(pkgDir, "/"+strings.TrimSuffix(targetDir, "/")+"/")
}

func matchesPattern(importPath string, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(importPath, prefix)
	}
	return importPath == pattern
}