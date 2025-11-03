# CI/CD Documentation

This document describes the CI/CD setup for the K4A project.

## Overview

The project uses GitHub Actions for continuous integration and deployment with the following workflows:

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:** Push and Pull Requests to `main`, `master`, or `develop` branches

**Jobs:**
- **Lint**: Runs `golangci-lint` to check code quality
- **Format Check**: Verifies code is properly formatted with `gofmt`
- **Test**: Runs unit tests with race detection and coverage reporting
- **Build**: Builds binaries for all supported platforms (Linux, macOS, Windows - amd64 & arm64)

### 2. Release Workflow (`.github/workflows/release.yml`)

**Triggers:** Push of tags matching `v*` (e.g., `v1.0.0`)

**Jobs:**
- **Changelog**: Generates changelog using git-cliff
- **Build Release**: Builds binaries for all platforms with version information
- **Docker**: Builds and pushes multi-arch Docker images to GitHub Container Registry
- **Release**: Creates GitHub release with binaries, checksums, and changelog

### 3. CodeQL Workflow (`.github/workflows/codeql.yml`)

**Triggers:** 
- Push and Pull Requests to `main`, `master`, or `develop` branches
- Weekly schedule (Mondays at midnight)

**Purpose:** Security scanning and code quality analysis

### 4. Auto Format Workflow (`.github/workflows/auto-format.yml`)

**Triggers:** Pull Requests to `main`, `master`, or `develop` branches

**Purpose:** Automatically formats code and commits changes if needed

## Dependabot Configuration

Dependabot is configured to automatically create PRs for:
- Go module updates (weekly)
- GitHub Actions updates (weekly)
- Docker base image updates (weekly)

## Creating a Release

To create a new release:

1. Ensure all changes are merged to `main`
2. Create and push a tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Generate changelog from conventional commits
   - Create a GitHub release
   - Publish Docker images

## Version Information

The build process injects version information into the binary:
- `version`: Git tag (e.g., `v1.0.0`)
- `commit`: Git commit SHA (short)
- `date`: Build timestamp

Users can check version with:
```bash
k4a -v
# or
k4a --version
```

## Conventional Commits

The project uses [Conventional Commits](https://www.conventionalcommits.org/) for changelog generation:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code formatting
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Test additions/updates
- `chore:` - Maintenance tasks

## Docker Images

Docker images are published to GitHub Container Registry (GHCR):

```bash
ghcr.io/smart-fellas/k4a:latest
ghcr.io/smart-fellas/k4a:v1.0.0
ghcr.io/smart-fellas/k4a:1.0
ghcr.io/smart-fellas/k4a:1
```

Multi-arch support: `linux/amd64`, `linux/arm64`

## Binary Naming Convention

Release binaries follow this pattern:
```
k4a-{os}-{arch}[.exe]
```

Examples:
- `k4a-linux-amd64`
- `k4a-darwin-arm64`
- `k4a-windows-amd64.exe`

## Checksums

Each release includes a `checksums.txt` file with SHA256 checksums of all binaries for verification.

## Local Testing

Test the CI locally before pushing:

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Build for all platforms
make build-all

# Check version info
./bin/k4a -v
```

## Troubleshooting

### CI Failures

1. **Lint failures**: Run `make lint` locally and fix issues
2. **Format failures**: Run `make fmt` locally
3. **Test failures**: Run `make test` locally
4. **Build failures**: Check Go version and dependencies

### Release Issues

1. **No release created**: Ensure tag matches `v*` pattern
2. **Missing binaries**: Check build job logs
3. **Docker push failed**: Verify GITHUB_TOKEN permissions

## Permissions Required

The workflows require the following GitHub token permissions:
- `contents: write` - For creating releases
- `packages: write` - For pushing Docker images
- `security-events: write` - For CodeQL scanning

These are automatically available via `GITHUB_TOKEN`.

