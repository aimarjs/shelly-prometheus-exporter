# Contributing to Shelly Prometheus Exporter

Thank you for your interest in contributing to the Shelly Prometheus Exporter! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Adding Device Support](#adding-device-support)
- [Documentation](#documentation)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/shelly-prometheus-exporter.git
   cd shelly-prometheus-exporter
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/aimarjs/shelly-prometheus-exporter.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Docker (optional, for containerized testing)
- Access to Shelly devices for testing

### Local Development

1. **Install dependencies**:

   ```bash
   go mod download
   ```

2. **Build the project**:

   ```bash
   make build
   ```

3. **Run tests**:

   ```bash
   make test
   ```

4. **Run linting**:
   ```bash
   make lint
   ```

### Configuration

Create a `.shelly-exporter.yaml` file for local testing:

```yaml
listen_address: ":8080"
shelly_devices:
  - "http://192.168.1.100" # Your Shelly device IP
log_level: "debug"
```

## Making Changes

### Branch Strategy

- Create a new branch for each feature or bugfix
- Use descriptive branch names:
  - `feature/add-shelly-plus-1pm`
  - `fix/relay-metrics-collection`
  - `docs/update-readme`

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` for formatting
- Add comments for exported functions and types
- Keep functions small and focused

### Commit Messages

Use clear, descriptive commit messages:

```
feat: add support for Shelly Plus 1PM device

- Implement RPC API client for Plus series
- Add relay and power monitoring metrics
- Update documentation with new device info

Closes #123
```

## Testing

### Unit Tests

Write unit tests for new functionality:

```bash
go test ./internal/...
```

### Integration Tests

Test with actual Shelly devices:

1. Configure your device in `.shelly-exporter.yaml`
2. Run the exporter: `./shelly-exporter --config=.shelly-exporter.yaml`
3. Check metrics: `curl http://localhost:8080/metrics`
4. Verify all expected metrics are present

### Test Coverage

Maintain good test coverage:

```bash
make test-coverage
```

## Submitting Changes

### Pull Request Process

1. **Update your fork**:

   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Create a pull request**:

   - Use the provided PR template
   - Include a clear description of changes
   - Reference related issues
   - Add screenshots if applicable

3. **Ensure CI passes**:
   - All tests must pass
   - Code must be properly formatted
   - No linting errors

### Review Process

- Maintainers will review your PR
- Address feedback promptly
- Keep PRs focused and reasonably sized
- Be responsive to questions and suggestions

## Adding Device Support

### Device Support Checklist

When adding support for a new Shelly device:

1. **Research the device**:

   - Check official documentation
   - Identify API endpoints
   - Understand device capabilities

2. **Implement the client**:

   - Add device-specific API calls
   - Handle different response formats
   - Implement error handling

3. **Add metrics collection**:

   - Define relevant metrics
   - Map device data to Prometheus format
   - Handle missing data gracefully

4. **Update documentation**:

   - Add device to supported list
   - Document specific metrics
   - Provide configuration examples

5. **Test thoroughly**:
   - Test with actual device
   - Verify all metrics work correctly
   - Test error scenarios

### Device Support Template

```go
// internal/client/device_shelly_plus_1pm.go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

// ShellyPlus1PMClient handles communication with Shelly Plus 1PM devices
type ShellyPlus1PMClient struct {
    *Client
}

// GetStatus retrieves status from Shelly Plus 1PM device
func (c *ShellyPlus1PMClient) GetStatus(ctx context.Context) (*StatusResponse, error) {
    // Implementation here
}
```

## Documentation

### Code Documentation

- Document all exported functions and types
- Include usage examples where helpful
- Keep documentation up to date

### User Documentation

- Update README.md for new features
- Add configuration examples
- Document breaking changes in CHANGELOG.md

### API Documentation

- Document new metrics in README.md
- Include metric descriptions and labels
- Provide example queries

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Changelog

Update `CHANGELOG.md` for each release:

- Group changes by type (Added, Changed, Fixed, Removed)
- Include issue references
- Follow the existing format

## Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and general discussion
- **Documentation**: Check the README and docs/ directory

## Recognition

Contributors will be recognized in:

- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to the Shelly Prometheus Exporter! 🎉
