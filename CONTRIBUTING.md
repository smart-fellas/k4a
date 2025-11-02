# Contributing to K4A

Thank you for your interest in contributing to K4A! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.21 or later
- [kafkactl](https://github.com/michelin/kafkactl) installed and configured
- Git

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/k4a.git
   cd k4a
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/smart-fellas/k4a.git
   ```
4. Install dependencies:
   ```bash
   make deps
   ```

## Development Workflow

### Building

```bash
make build
```

### Running

```bash
make run
```

### Testing

```bash
make test
```

### Formatting

We use `gofmt` for code formatting. Before committing, format your code:

```bash
make fmt
```

This will automatically format all Go files in the project.

### Linting

We use `golangci-lint` for linting. Run it locally before pushing:

```bash
make lint
```

Fix any issues reported by the linter.

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

- `feat: add new feature` - New features
- `fix: resolve bug in component` - Bug fixes
- `docs: update documentation` - Documentation changes
- `style: format code` - Code style changes (formatting, etc.)
- `refactor: restructure code` - Code refactoring
- `perf: improve performance` - Performance improvements
- `test: add tests` - Adding or updating tests
- `chore: update dependencies` - Maintenance tasks

### Examples

```
feat: add pause/resume functionality for connectors
fix: resolve crash when listing empty topics
docs: update installation instructions
refactor: simplify topic list rendering
test: add unit tests for config loader
```

## Pull Request Process

1. Create a feature branch:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

3. Push to your fork:
   ```bash
   git push origin feat/your-feature-name
   ```

4. Open a Pull Request against the `main` branch

5. Ensure all CI checks pass:
   - Code formatting
   - Linting
   - Tests
   - Build for all platforms

6. Wait for review and address any feedback

## Code Style

- Follow standard Go conventions
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Use early returns to reduce nesting

## Testing

- Write tests for new features
- Update tests when modifying existing code
- Ensure tests pass locally before pushing
- Aim for good test coverage

## Project Structure

```
k4a/
├── cmd/k4a/          # Main application entry point
├── internal/         # Internal packages
│   ├── app/          # Application logic
│   ├── config/       # Configuration handling
│   ├── kafkactl/     # Kafkactl client wrapper
│   └── ui/           # UI components and views
├── pkg/              # Public packages
└── .github/          # GitHub workflows and configs
```

## Release Process

Releases are automated via GitHub Actions:

1. Ensure all changes are merged to `main`
2. Create and push a tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Generate changelog
   - Create GitHub release
   - Build and push Docker image

## Getting Help

- Open an issue for bug reports or feature requests
- Join discussions in GitHub Discussions
- Check existing issues and PRs before creating new ones

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Focus on what is best for the community

Thank you for contributing to K4A! 🎉

