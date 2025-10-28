# Contributing to CoreKV

Thank you for considering contributing to CoreKV! This document outlines the process and guidelines for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR-USERNAME/corekv.git`
3. Create a new branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `make test`
6. Run linters: `make lint`
7. Commit your changes: `git commit -m "Add my feature"`
8. Push to your fork: `git push origin feature/my-feature`
9. Create a Pull Request

## Development Guidelines

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` and `goimports` for formatting (run `make fmt`)
- Keep functions small and focused
- Add comments for exported functions and types
- Write clear commit messages

### Testing

- Write tests for new features
- Ensure all tests pass before submitting a PR
- Aim for good test coverage
- Run `make test` to execute all tests

### Linting

- Run `make lint` before committing
- Address all linter warnings and errors
- Pre-commit hooks will run automatically if installed

### Pre-commit Hooks

Install pre-commit hooks to catch issues before committing:

```bash
pip install pre-commit
make install-hooks
```

## Pull Request Process

1. Update the README.md with details of changes if applicable
2. Ensure all tests pass and linters are happy
3. Update documentation as needed
4. The PR will be merged once reviewed and approved

## Code Review

All contributions require code review. We use GitHub pull requests for this purpose.

## Questions?

Feel free to open an issue for any questions or concerns!
