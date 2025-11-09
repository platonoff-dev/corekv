# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Using Task (preferred)
- `task build` - Build the project
- `task test` - Run tests with race detection and coverage
- `task test-coverage` - Run tests and generate HTML coverage report
- `task lint` - Run golangci-lint
- `task fmt` - Format code with go fmt and goimports
- `task vet` - Run go vet
- `task clean` - Clean build artifacts and coverage files
- `task all` - Run formatting, linting, tests, and build

### Alternative commands
If Task is not available, use these Go commands directly:
- `go build -v ./...` - Build the project
- `go test -v -race -coverprofile=coverage.out ./...` - Run tests
- `golangci-lint run ./...` - Run linting
- `go fmt ./...` - Format code
- `go vet ./...` - Run go vet

### Development setup
- `task install-tools` - Install golangci-lint and goimports
- `task install-hooks` - Install pre-commit hooks
- `task mod-tidy` - Tidy go.mod

## Architecture

CoreKV is a persistent key-value database implementation using the Bitcask storage engine pattern.

### Storage Layer
The core storage implementation is in `storage/bitcask/`:

- **Engine** (`engine.go`): Main storage engine with in-memory index and file management
- **DataEntry** (`entry.go`): Binary format for log entries with CRC32 checksums
- **Index** (`index.go`): In-memory hash table mapping keys to file positions
- **File format**: Append-only log files with entries containing timestamp, key/value lengths, and data

### Key Components
- Uses LittleEndian binary encoding for all data structures
- CRC32 checksums with Castagnoli polynomial for data integrity
- Thread-safe operations (as indicated in README)
- Time-based timestamps for entries

### Current Implementation Status
The codebase appears to be in active development:
- Basic Bitcask engine structure is implemented
- Entry encoding is implemented but decoding is incomplete (`entry.go:49`)
- No file rotation or merge operations yet (mentioned in git history)

## Testing
- Tests use testify framework
- Run with race detection enabled by default
- Coverage reports generated in `coverage.out` and `coverage.html`

## Linting
Comprehensive linting setup with golangci-lint including security, performance, and style checks. Configuration in `.golangci.yml` with 50+ enabled linters.
