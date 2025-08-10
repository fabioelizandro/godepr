# godepr

A Go dependency governance tool that enables programmers to enforce rules on package dependencies within their projects.

## Overview

`godepr` is a lightweight command-line tool that helps maintain clean dependency boundaries in Go projects. It reads configuration from a `.godepr` file and validates that imports across your codebase comply with defined rules.

## Installation

```bash
go install github.com/yourusername/godepr@latest
```

## Usage

```bash
godepr entrypoint.go
```

The tool will read the `.godepr` configuration file from the current directory and validate all imports starting from the specified entrypoint.

## Configuration

Create a `.godepr` file in your project root:

```json
{
    "rules": [
        {
            "directory": "internal/domain/",
            "ruletype": "denied-list",
            "rulebody": ["github.com/external/*", "internal/infrastructure/*"]
        },
        {
            "directory": "internal/api/",
            "ruletype": "allowed-list",
            "rulebody": ["internal/domain/*", "internal/services/*"]
        }
    ]
}
```

## Rule Types

### `denied-list`
Starts open and blocks specific dependencies matching the patterns.
- **directory**: Directory path to apply the rule
- **rulebody**: Array of regex patterns for denied imports

### `allowed-list` (planned)
Starts closed and only allows specific dependencies matching the patterns.

### `std-dep-only` (planned)
Only allows standard library dependencies with optional exceptions.

### `deny-sibling-deps` (planned)
Prevents packages from importing sibling packages within the same directory.

### `name-convention` (planned)
Enforces naming conventions for packages within the directory.

## Development

### Project Structure

```
godepr/
├── .godepr          # Configuration file
├── .gitignore       # Git ignore rules
├── CLAUDE.md        # Development notes
├── README.md        # This file
├── go.mod           # Go module definition
├── main.go          # Main entry point
├── godepr.go        # Core logic
└── godepr_test.go   # Tests
```

### Testing

Run tests with:
```bash
go test -v ./...
```

### Git Hooks

Add to `.git/hooks/pre-commit`:
```bash
#!/bin/sh
godepr main.go
```

## Current Implementation Status

- [x] Project setup
- [x] Basic structure
- [ ] `.godepr` file parsing
- [ ] Import extraction from Go files
- [ ] Rule validation engine
- [ ] denied-list rule implementation
- [ ] allowed-list rule implementation
- [ ] std-dep-only rule implementation
- [ ] deny-sibling-deps rule implementation
- [ ] name-convention rule implementation
- [ ] CLI interface
- [ ] Error reporting
- [ ] Performance optimization

## Contributing

This is an open-source project. Contributions are welcome!

## License

MIT