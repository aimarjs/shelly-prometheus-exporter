## [1.1.1](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.1.0...v1.1.1) (2025-10-14)


### Bug Fixes

* address linter warnings in network retry implementation ([2cc7d6d](https://github.com/aimarjs/shelly-prometheus-exporter/commit/2cc7d6db627ad38c120760a40f5abc49fb31bcf6))
* change device status card color mode to background ([412362b](https://github.com/aimarjs/shelly-prometheus-exporter/commit/412362b1e8bd53ee756b77a04b439913528cf5af))
* create proper Grafana provisioning format ([e51dec5](https://github.com/aimarjs/shelly-prometheus-exporter/commit/e51dec592c17ae47d38887ffdfeee3dcc3b9ec36))
* create working Grafana dashboard ([6a32f8b](https://github.com/aimarjs/shelly-prometheus-exporter/commit/6a32f8bc1ab73fae685d3aceaf39e685c0c3351c))
* format heating-cost-dashboard.json for consistency ([80120d1](https://github.com/aimarjs/shelly-prometheus-exporter/commit/80120d1ca108c9b97a104f7cc3bae22c06cf643d))
* remove dashboard wrapper from JSON files for Grafana provisioning ([ac6a1dc](https://github.com/aimarjs/shelly-prometheus-exporter/commit/ac6a1dcfef56b045cddad1b6c32e4a2f447824ff))
* replace increase() with rate() for energy metrics ([edc38be](https://github.com/aimarjs/shelly-prometheus-exporter/commit/edc38be9f7c6396417c813fe26134e3b61f4953e))
* resolve persistent network connectivity issues with connection pooling and retry logic ([7f10b53](https://github.com/aimarjs/shelly-prometheus-exporter/commit/7f10b53820e40c570fd9afdb714b65cf70481fd1))
* update heating-cost-dashboard.json and shelly-dashboard.json for improved structure ([d4160af](https://github.com/aimarjs/shelly-prometheus-exporter/commit/d4160af8e1e1bfe4baa23919b202fa46397000b0))
* update Prometheus rules to match current metric structure ([c631c84](https://github.com/aimarjs/shelly-prometheus-exporter/commit/c631c84936d10666741864662a00bb1c5a2fcab9))
* update shelly-simple-dashboard.json for improved structure ([5243a71](https://github.com/aimarjs/shelly-prometheus-exporter/commit/5243a71e3481a165af72ae65d03c8be4f600e2af))
* update title in shelly-dashboard.json for clarity ([fb7f864](https://github.com/aimarjs/shelly-prometheus-exporter/commit/fb7f864c981cddb40d0f0eecb62866405d6331a4))

# [1.1.0](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.12...v1.1.0) (2025-10-06)


### Bug Fixes

* add fallback values to Prometheus recording rules ([ad1e566](https://github.com/aimarjs/shelly-prometheus-exporter/commit/ad1e56671816ed29d1f58176e02eba1e34bc4e2d))
* add missing newline at end of mocks.go file ([64633da](https://github.com/aimarjs/shelly-prometheus-exporter/commit/64633da5baf7f193522708e5aa2b756f9a5b5f41))
* configure Qlty to use Go 1.24.0 toolchain ([c6b7a35](https://github.com/aimarjs/shelly-prometheus-exporter/commit/c6b7a35b5e3a9198d196470a442fc0d23f7d126e))
* configure Qlty to use Go 1.25.0 ([667feb2](https://github.com/aimarjs/shelly-prometheus-exporter/commit/667feb26860ef465e28ab7ec03a629d63bda7b21))
* configure Qlty to work with Go 1.21 ([13cd3a5](https://github.com/aimarjs/shelly-prometheus-exporter/commit/13cd3a5b5edbcba8fadc65d39d5e69123238516a))
* correct Grafana datasource provisioning format ([e6abdfe](https://github.com/aimarjs/shelly-prometheus-exporter/commit/e6abdfe8a71c533247063d2513057e67d7953c43))
* correct type mismatch in createLegacyResponse function ([c0e93a8](https://github.com/aimarjs/shelly-prometheus-exporter/commit/c0e93a8ceaae98bfee7755d457fccfbc752037dc))
* downgrade Go version from 1.25 to 1.24 for Qlty compatibility ([27a0959](https://github.com/aimarjs/shelly-prometheus-exporter/commit/27a0959914224a07be0b2ae1efa1e2a485fa33b2))
* downgrade Go version to 1.23 across configurations ([4d369cf](https://github.com/aimarjs/shelly-prometheus-exporter/commit/4d369cf0e9dababcba5c6eb178c4488753ff474c))
* downgrade Go version to 1.24 for Qlty compatibility ([ec0cd02](https://github.com/aimarjs/shelly-prometheus-exporter/commit/ec0cd0215b7bafff145acd784289721ed9384de6))
* final Qlty check fixes ([7581f68](https://github.com/aimarjs/shelly-prometheus-exporter/commit/7581f688bba1559015d3956f392b515cb5a81c37))
* handle equal start and end times in GetCurrentRate() ([2a9c027](https://github.com/aimarjs/shelly-prometheus-exporter/commit/2a9c02738e7a8e34d9175e1d46a172dcf14ac4f2))
* handle time ranges that cross midnight ([3c6d9c5](https://github.com/aimarjs/shelly-prometheus-exporter/commit/3c6d9c58a39b3fad8ac516a6d68ab1de1b34ef9b))
* make TestCollector_Describe more flexible and robust ([efb7bd7](https://github.com/aimarjs/shelly-prometheus-exporter/commit/efb7bd7793ad33a6afe344fb007fa2a2ab82d195))
* remove unnecessary whitespace in metrics.go ([33ffc85](https://github.com/aimarjs/shelly-prometheus-exporter/commit/33ffc85f62df523396c14ad0d06f0a0d5c07d54a))
* remove unused variable in metrics test helper function ([09f2866](https://github.com/aimarjs/shelly-prometheus-exporter/commit/09f2866110ad191ed7c0742ce201d52a60ce8ac7))
* resolve all Qlty check issues ([de1b285](https://github.com/aimarjs/shelly-prometheus-exporter/commit/de1b28562f984c1a4782951334e10b4de3a10f8a))
* resolve deadlock in heating percentage calculation ([3cdef1b](https://github.com/aimarjs/shelly-prometheus-exporter/commit/3cdef1b878be36b959063c69a77594af4264a911))
* resolve Qlty configuration and code quality issues ([71c47e9](https://github.com/aimarjs/shelly-prometheus-exporter/commit/71c47e91ba69cddede94c1d329db74e501123dfc))
* restore Grafana datasource provisioning configuration ([d383895](https://github.com/aimarjs/shelly-prometheus-exporter/commit/d38389521621e130f2e41ac1642399b03454a96f))
* revert undefined HeatingPercentageTimeout field reference ([0be1ab4](https://github.com/aimarjs/shelly-prometheus-exporter/commit/0be1ab4b37fd4ca1669e962d41e64f0f0cb9fe21))
* update Go version to 1.25 across all configuration files ([71bb5cb](https://github.com/aimarjs/shelly-prometheus-exporter/commit/71bb5cb9b04b98ad8d85c01bf3a33a844d619fcf))
* update golangci-lint to version 2.5.0 for Go 1.25 compatibility ([9ebf088](https://github.com/aimarjs/shelly-prometheus-exporter/commit/9ebf08856b64cd9a99e30ad520c2a4fa0d7ec492))


### Features

* implement Phase 1 configuration enhancements ([8f912e4](https://github.com/aimarjs/shelly-prometheus-exporter/commit/8f912e4db8b68d038dd68821b6f762ac37be6e02))
* implement Phase 2 metric enhancements ([95129e1](https://github.com/aimarjs/shelly-prometheus-exporter/commit/95129e12ba3824c01a636059dab03f96edf22a02))
* implement Phase 3 dashboard integration ([04212fc](https://github.com/aimarjs/shelly-prometheus-exporter/commit/04212fc9985ce4c7fcf3f2fca039c603f2ca8139))
* implement time-based rate calculation ([149680e](https://github.com/aimarjs/shelly-prometheus-exporter/commit/149680e07ed3709ad8fca92f590fec1a0352911c))

## [1.0.12](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.11...v1.0.12) (2025-10-04)


### Bug Fixes

* use GITHUB_TOKEN for Docker registry authentication ([03eda9a](https://github.com/aimarjs/shelly-prometheus-exporter/commit/03eda9a09046c80057330b40c08e88921e984d55))

## [1.0.11](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.10...v1.0.11) (2025-10-04)


### Bug Fixes

* remove buildx from dockers configuration ([0422463](https://github.com/aimarjs/shelly-prometheus-exporter/commit/0422463e7d01ad25a757ed1a0b0177b310aca39f))

## [1.0.10](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.9...v1.0.10) (2025-10-04)


### Bug Fixes

* revert to stable dockers configuration ([fe9db04](https://github.com/aimarjs/shelly-prometheus-exporter/commit/fe9db04cd25cc888f074daee6e475e381ef27dc3))

## [1.0.9](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.8...v1.0.9) (2025-10-04)


### Bug Fixes

* use platform-specific binary path in Dockerfile ([1433246](https://github.com/aimarjs/shelly-prometheus-exporter/commit/14332462594806b700a69201c2f29d8f39bedef1))

## [1.0.8](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.7...v1.0.8) (2025-10-04)


### Bug Fixes

* remove config file copy from Dockerfile ([8806a71](https://github.com/aimarjs/shelly-prometheus-exporter/commit/8806a716e64b7cd362d2c969a8f4490e219a1af3))

## [1.0.7](https://github.com/aimarjs/shelly-prometheus-exporter/compare/v1.0.6...v1.0.7) (2025-10-03)


### Bug Fixes

* simplify Dockerfile for GoReleaser dockers_v2 ([64b93d7](https://github.com/aimarjs/shelly-prometheus-exporter/commit/64b93d7a0e096cae7e1d79a1b7776f2646bc3336))
* update GoReleaser configuration for archive formats and version template ([86ade10](https://github.com/aimarjs/shelly-prometheus-exporter/commit/86ade10eeb89b6c1c65e6c2f3c498917f23b3c81))

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
