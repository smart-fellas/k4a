# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GitHub Actions CI/CD pipeline
- Multi-platform binary releases (Linux, macOS, Windows - amd64 & arm64)
- Docker image support with GHCR
- Automated changelog generation
- CodeQL security scanning
- Dependabot for dependency updates
- Code formatting and linting in CI
- Version flag support (`k4a -v` or `k4a --version`)

### Changed
- Enhanced Makefile with version injection and multi-platform builds

## [0.1.0] - Initial Release

### Added
- Terminal UI for Kafka management
- Topics view and management
- Schemas view and management
- Connectors view and management
- k9s-like navigation
- Multi-context support
- YAML resource descriptions
- Real-time updates

[Unreleased]: https://github.com/smart-fellas/k4a/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/smart-fellas/k4a/releases/tag/v0.1.0

