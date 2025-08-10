package main

import (
	"testing"
)

func TestCheckDeniedListRule(t *testing.T) {
	t.Run("blocks imports matching denied patterns", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports: []string{
					"fmt",
					"github.com/blocked/package",
					"internal/allowed/package",
				},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 1 {
			t.Fatalf("expected 1 violation, got %d", len(violations))
		}

		if violations[0].ImportPath != "github.com/blocked/package" {
			t.Errorf("expected violation for 'github.com/blocked/package', got %s", violations[0].ImportPath)
		}
	})

	t.Run("allows imports not matching denied patterns", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports: []string{
					"fmt",
					"internal/allowed/package",
					"github.com/allowed/package",
				},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 0 {
			t.Fatalf("expected 0 violations, got %d", len(violations))
		}
	})

	t.Run("handles multiple denied patterns", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports: []string{
					"fmt",
					"github.com/blocked/package",
					"internal/infrastructure/db",
					"internal/domain/model",
				},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*", "internal/infrastructure/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 2 {
			t.Fatalf("expected 2 violations, got %d", len(violations))
		}

		expectedViolations := map[string]bool{
			"github.com/blocked/package": false,
			"internal/infrastructure/db":  false,
		}

		for _, v := range violations {
			if _, exists := expectedViolations[v.ImportPath]; exists {
				expectedViolations[v.ImportPath] = true
			}
		}

		for path, found := range expectedViolations {
			if !found {
				t.Errorf("expected violation for %s not found", path)
			}
		}
	})

	t.Run("only applies to files in specified directory", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/other",
				Dir:        "/project/other/",
				Imports: []string{
					"github.com/blocked/package",
				},
			},
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports:    []string{},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 0 {
			t.Fatalf("expected 0 violations for file outside rule directory, got %d", len(violations))
		}
	})

	t.Run("handles exact match patterns", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports: []string{
					"github.com/exact/match",
					"github.com/exact/other",
				},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/exact/match"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 1 {
			t.Fatalf("expected 1 violation, got %d", len(violations))
		}

		if violations[0].ImportPath != "github.com/exact/match" {
			t.Errorf("expected violation for 'github.com/exact/match', got %s", violations[0].ImportPath)
		}
	})

	t.Run("handles empty imports gracefully", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports:    []string{},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 0 {
			t.Fatalf("expected 0 violations for file with no imports, got %d", len(violations))
		}
	})

	t.Run("handles named imports", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports: []string{
					"github.com/blocked/package",
					"github.com/blocked/sideeffect",
					"github.com/blocked/dot",
				},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 3 {
			t.Fatalf("expected 3 violations, got %d", len(violations))
		}
	})

	t.Run("handles subdirectories correctly", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/example/sub",
				Dir:        "/project/example/sub/",
				Imports: []string{
					"github.com/blocked/package",
				},
			},
			{
				ImportPath: "testproject/example",
				Dir:        "/project/example/",
				Imports: []string{
					"github.com/blocked/other",
				},
			},
		}

		rule := Rule{
			Directory: "example/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 2 {
			t.Fatalf("expected 2 violations (one from each package), got %d", len(violations))
		}
	})

	t.Run("handles packages with no directory match", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/pkg",
				Dir:        "/project/pkg/",
				Imports: []string{
					"github.com/blocked/package",
				},
			},
		}

		rule := Rule{
			Directory: "nonexistent/",
			RuleType:  "denied-list",
			RuleBody:  []string{"github.com/blocked/*"},
		}

		violations := CheckDeniedListRule(packages, rule)

		if len(violations) != 0 {
			t.Fatalf("expected 0 violations for non-matching directory, got %d", len(violations))
		}
	})
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		pattern    string
		want       bool
	}{
		{
			name:       "exact match",
			importPath: "github.com/user/package",
			pattern:    "github.com/user/package",
			want:       true,
		},
		{
			name:       "wildcard at end matches",
			importPath: "github.com/user/package/subpackage",
			pattern:    "github.com/user/package/*",
			want:       true,
		},
		{
			name:       "wildcard at end matches exact",
			importPath: "github.com/user/package",
			pattern:    "github.com/user/package*",
			want:       true,
		},
		{
			name:       "wildcard does not match different prefix",
			importPath: "github.com/other/package",
			pattern:    "github.com/user/package/*",
			want:       false,
		},
		{
			name:       "no match without wildcard",
			importPath: "github.com/user/package/sub",
			pattern:    "github.com/user/package",
			want:       false,
		},
		{
			name:       "wildcard without slash matches",
			importPath: "github.com/user/package",
			pattern:    "github.com/user/package*",
			want:       true,
		},
		{
			name:       "wildcard without slash matches subpackage",
			importPath: "github.com/user/package/sub",
			pattern:    "github.com/user/package*",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesPattern(tt.importPath, tt.pattern)
			if got != tt.want {
				t.Errorf("matchesPattern(%q, %q) = %v, want %v", tt.importPath, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestCheckRules(t *testing.T) {
	t.Run("processes multiple rules", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/domain",
				Dir:        "/project/domain/",
				Imports: []string{
					"github.com/external/lib",
				},
			},
			{
				ImportPath: "testproject/api",
				Dir:        "/project/api/",
				Imports: []string{
					"github.com/gin-gonic/gin",
				},
			},
		}

		rules := []Rule{
			{
				Directory: "domain/",
				RuleType:  "denied-list",
				RuleBody:  []string{"github.com/external/*"},
			},
			{
				Directory: "api/",
				RuleType:  "denied-list",
				RuleBody:  []string{"github.com/gin-gonic/*"},
			},
		}

		violations := CheckRules(packages, rules)

		if len(violations) != 2 {
			t.Fatalf("expected 2 violations, got %d", len(violations))
		}
	})

	t.Run("handles unsupported rule types", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject",
				Dir:        "/project/",
				Imports:    []string{},
			},
		}

		rules := []Rule{
			{
				Directory: "./",
				RuleType:  "unsupported-type",
				RuleBody:  []string{"some-pattern"},
			},
		}

		violations := CheckRules(packages, rules)

		if len(violations) != 0 {
			t.Fatalf("expected 0 violations for unsupported rule type, got %d", len(violations))
		}
	})

	t.Run("applies multiple rules to same package", func(t *testing.T) {
		packages := []Package{
			{
				ImportPath: "testproject/api",
				Dir:        "/project/api/",
				Imports: []string{
					"github.com/gin-gonic/gin",
					"github.com/external/lib",
				},
			},
		}

		rules := []Rule{
			{
				Directory: "api/",
				RuleType:  "denied-list",
				RuleBody:  []string{"github.com/gin-gonic/*"},
			},
			{
				Directory: "api/",
				RuleType:  "denied-list",
				RuleBody:  []string{"github.com/external/*"},
			},
		}

		violations := CheckRules(packages, rules)

		if len(violations) != 2 {
			t.Fatalf("expected 2 violations, got %d", len(violations))
		}
	})
}

func TestParsePackagesJSON(t *testing.T) {
	t.Run("parses single package JSON", func(t *testing.T) {
		jsonData := `{"ImportPath":"testproject/example","Dir":"/project/example","Imports":["fmt","strings"]}`

		packages, err := ParsePackagesJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(packages) != 1 {
			t.Fatalf("expected 1 package, got %d", len(packages))
		}

		if packages[0].ImportPath != "testproject/example" {
			t.Errorf("expected ImportPath 'testproject/example', got %s", packages[0].ImportPath)
		}

		if len(packages[0].Imports) != 2 {
			t.Errorf("expected 2 imports, got %d", len(packages[0].Imports))
		}
	})

	t.Run("parses multiple packages JSON", func(t *testing.T) {
		jsonData := `{"ImportPath":"testproject/pkg1","Dir":"/project/pkg1","Imports":["fmt"]}
{"ImportPath":"testproject/pkg2","Dir":"/project/pkg2","Imports":["strings"]}`

		packages, err := ParsePackagesJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(packages) != 2 {
			t.Fatalf("expected 2 packages, got %d", len(packages))
		}

		if packages[0].ImportPath != "testproject/pkg1" {
			t.Errorf("expected first package ImportPath 'testproject/pkg1', got %s", packages[0].ImportPath)
		}

		if packages[1].ImportPath != "testproject/pkg2" {
			t.Errorf("expected second package ImportPath 'testproject/pkg2', got %s", packages[1].ImportPath)
		}
	})

	t.Run("handles empty JSON", func(t *testing.T) {
		jsonData := ``

		packages, err := ParsePackagesJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(packages) != 0 {
			t.Fatalf("expected 0 packages, got %d", len(packages))
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		jsonData := `{invalid json}`

		_, err := ParsePackagesJSON(jsonData)
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})
}

func TestIsInDirectory(t *testing.T) {
	tests := []struct {
		name      string
		pkgDir    string
		targetDir string
		want      bool
	}{
		{
			name:      "exact match with trailing slashes",
			pkgDir:    "/project/example/",
			targetDir: "example/",
			want:      true,
		},
		{
			name:      "exact match without trailing slashes",
			pkgDir:    "/project/example",
			targetDir: "example/",
			want:      true,
		},
		{
			name:      "subdirectory match",
			pkgDir:    "/project/example/sub/",
			targetDir: "example/",
			want:      true,
		},
		{
			name:      "no match - different directory",
			pkgDir:    "/project/other/",
			targetDir: "example/",
			want:      false,
		},
		{
			name:      "no match - parent directory",
			pkgDir:    "/project/",
			targetDir: "example/",
			want:      false,
		},
		{
			name:      "absolute path match",
			pkgDir:    "/home/user/project/api/",
			targetDir: "api/",
			want:      true,
		},
		{
			name:      "nested subdirectory match",
			pkgDir:    "/project/api/v1/handlers/",
			targetDir: "api/",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isInDirectory(tt.pkgDir, tt.targetDir)
			if got != tt.want {
				t.Errorf("isInDirectory(%q, %q) = %v, want %v", tt.pkgDir, tt.targetDir, got, tt.want)
			}
		})
	}
}