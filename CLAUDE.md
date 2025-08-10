# CLAUDE.md - Development Notes for godepr

## Project Overview
`godepr` is a Go dependency governance tool that validates import rules defined in a `.godepr` configuration file.

## Testing Commands
```bash
go test -v ./...
go test -race ./...
go test -cover ./...
```

## Linting and Type Checking
```bash
go vet ./...
gofmt -s -w .
```

## Build Commands
```bash
go build -o godepr main.go
```

## Key Design Decisions

### Testing Philosophy
- Following the testing style from https://github.com/fabioelizandro/testfill/blob/main/testfill_test.go
- Tests should be clear, descriptive, and focus on behavior
- Use table-driven tests where appropriate
- Each test should have a clear name describing what it tests

### Implementation Approach
1. Start with the simplest rule: `denied-list`
2. Build incrementally, test-first development
3. Parse actual Go AST instead of regex for accuracy
4. Keep the tool fast and suitable for git hooks

### Architecture Notes
- Core logic separated from CLI interface
- Rule validation is pluggable (strategy pattern)
- Each rule type implements a common interface
- Configuration parsing is separate from rule execution

## Current Focus
Implementing the `denied-list` rule with hardcoded configuration first, then adding JSON parsing.

## TODO for Next Sessions
- [ ] Add JSON configuration parsing
- [ ] Implement remaining rule types
- [ ] Add recursive directory walking
- [ ] Add caching for performance
- [ ] Add detailed error reporting with file:line references
- [ ] Add --fix flag for auto-fixing certain violations
- [ ] Add configuration validation
- [ ] Add support for .godepr.yaml as alternative format

## Performance Considerations
- Use parallel processing for large codebases
- Cache parsed ASTs when checking multiple rules
- Early exit on first violation (unless --all flag is used)

## Error Handling Strategy
- Clear error messages with file path and line number
- Suggest fixes where possible
- Exit with non-zero code on violations (for git hooks)

## Testing Checklist
- [ ] Unit tests for each rule type
- [ ] Integration tests with sample Go projects
- [ ] Edge cases (empty files, no imports, circular deps)
- [ ] Performance tests with large codebases
- [ ] Test with different Go versions