# CoreKV

Simple small persistent key-value database written in Go.

## Features

- Lightweight and fast persistent key-value storage
- Simple API
- Thread-safe operations

## Installation

```bash
go get github.com/platonoff-dev/corekv
```

## Usage

```bash
go run cmd/corekv/main.go
```

## Development

### Prerequisites

- Go 1.25 or higher
- golangci-lint (for linting)
- pre-commit (for pre-commit hooks)
- Task (https://taskfile.dev/) - optional but recommended (used via Taskfile.yml)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/platonoff-dev/corekv.git
cd corekv
```

2. Install development tools:
```bash
# Prefer Task if installed
if command -v task >/dev/null 2>&1; then
  task install-tools
else
  make install-tools
fi
```

3. Install pre-commit hooks:
```bash
pip install pre-commit  # if not already installed
# Prefer Task if installed
if command -v task >/dev/null 2>&1; then
  task install-hooks
else
  make install-hooks
fi
```

### Available development tasks

This project provides a `Taskfile.yml` (Task) and a `Makefile` for backward compatibility. Prefer `task` when available.

- `task build` - Build the project
- `task test` - Run tests
- `task test-coverage` - Run tests with coverage report
- `task lint` - Run linters
- `task fmt` - Format code
- `task clean` - Clean build artifacts
- `task all` - Run all checks and build

## Project Structure

```
.
├── cmd/
│   └── corekv/         # Main application
├── pkg/                # Public libraries
├── internal/           # Private application code
├── .github/
│   └── workflows/      # GitHub Actions workflows
└── Makefile            # Build and development tasks
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
