## [1.0.6](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.5...v1.0.6) (2025-10-03)


### Bug Fixes

* update dockers_v2 to use images and tags properties ([c7189b6](https://github.com/aimarjs/shelly-prometheus-exporter/commit/c7189b676a101537bf803c90a18c6ee126183652))

## [1.0.5](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.4...v1.0.5) (2025-10-03)


### Bug Fixes

* use correct image_templates property in dockers_v2 ([625e556](https://github.com/aimarjs/shelly-prometheus-exporter/commit/625e556da927132f37296edb09531100f23c76d6))

## [1.0.4](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.3...v1.0.4) (2025-10-03)


### Bug Fixes

* specify main package path in GoReleaser configuration ([ec81d97](https://github.com/aimarjs/shelly-prometheus-exporter/commit/ec81d97455d820fff2cce088428b8a8835ab0ebe))

## [1.0.3](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.2...v1.0.3) (2025-10-03)


### Bug Fixes

* checkout correct commit/tag in release job [skip-ci] ([524c700](https://github.com/aimarjs/shelly-prometheus-exporter/commit/524c7009d5ad1c7cbe17da5a4ab413a86a444be6))

## [1.0.2](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.1...v1.0.2) (2025-10-03)


### Bug Fixes

* allow GoReleaser to run on workflow_dispatch when semantic-release creates tags ([5a6e521](https://github.com/aimarjs/shelly-prometheus-exporter/commit/5a6e521ce9454d47d5e776ab67230efae988149b))

## [1.0.1](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.0...v1.0.1) (2025-10-03)


### Bug Fixes

* add Docker authentication to release workflow ([a40c71f](https://github.com/aimarjs/shelly-prometheus-exporter/commit/a40c71fcd7436ac159909c9915ac397c9084e863))
* resolve GoReleaser configuration and workflow issues ([c5a80ef](https://github.com/aimarjs/shelly-prometheus-exporter/commit/c5a80ef9090d2955622e2a8a227d654e35b1cb0d))

# 1.0.0 (2025-10-03)


### Bug Fixes

* add required permissions for semantic-release ([a344afd](https://github.com/aimarjs/shelly-prometheus-exporter/commit/a344afd0c98f4f6b7b0f686a5b3ba4a9bbbbe1a2))
* correct repository URLs from aimar to aimarjs ([a248266](https://github.com/aimarjs/shelly-prometheus-exporter/commit/a24826618db3d2a365da3428d742db6749be708f))


### Features

* add automated semantic versioning system ([41b9e96](https://github.com/aimarjs/shelly-prometheus-exporter/commit/41b9e96bddfc87154d979f5073596367876c4247))
* Enhance Shelly device support and metrics collection ([866de0a](https://github.com/aimarjs/shelly-prometheus-exporter/commit/866de0aac002db2593532d77c4ff60bae202f963))
* Extend support for Shelly Plug S and update documentation ([4c0faa0](https://github.com/aimarjs/shelly-prometheus-exporter/commit/4c0faa06ac2e1b143796459b669c74e18bc3459b))
* Update dependencies and enhance Shelly device metrics handling ([d42ffbd](https://github.com/aimarjs/shelly-prometheus-exporter/commit/d42ffbda08542890b708ba9b0bffd29136021bdf))
* update release workflow to use RELEASE_TOKEN for protected branch compatibility ([7095c59](https://github.com/aimarjs/shelly-prometheus-exporter/commit/7095c59ccc286ef797f157aa571a43750aee0604))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Shelly 1PM support with legacy API fallback
- Shelly Plug S support with legacy API
- Automatic device type detection
- Relay monitoring for Shelly 1PM and Plug S devices
- Unified metrics collection for multiple device types
- Initial project structure
- Basic Shelly device client
- Prometheus metrics collection
- HTTP server with health checks
- Configuration management
- Docker support
- Kubernetes manifests
- Comprehensive documentation

### Changed

- Updated client to support both RPC and legacy APIs
- Enhanced metrics collection for different device capabilities
- Extended legacy API support for Plug S devices

### Deprecated

### Removed

### Fixed

### Security

## [0.1.0] - 2024-01-XX

### Added

- Initial release
- Support for Shelly Pro3em and similar devices
- Basic metrics collection (power, relays, WiFi, temperature)
- TLS support for secure connections
- Docker and Kubernetes deployment options
- Configuration file support
- Command-line interface
