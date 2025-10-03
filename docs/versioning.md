# Versioning Guide

This project uses automated semantic versioning with [semantic-release](https://semantic-release.gitbook.io/) and [GoReleaser](https://goreleaser.com/).

## How It Works

### Automatic Versioning

- **Commit Message Convention**: Follow [Conventional Commits](https://www.conventionalcommits.org/) specification
- **Semantic Release**: Automatically determines version bumps based on commit messages
- **GoReleaser**: Builds and publishes releases with proper versioning

### Version Types

- **Major** (`1.0.0` → `2.0.0`): Breaking changes
- **Minor** (`1.0.0` → `1.1.0`): New features (backward compatible)
- **Patch** (`1.0.0` → `1.0.1`): Bug fixes (backward compatible)
- **Prerelease** (`1.0.0` → `1.0.1-beta.1`): Pre-release versions

## Commit Message Format

Use the following format for commit messages:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: New feature (triggers minor version bump)
- `fix`: Bug fix (triggers patch version bump)
- `docs`: Documentation changes (no version bump)
- `style`: Code style changes (no version bump)
- `refactor`: Code refactoring (no version bump)
- `test`: Test changes (no version bump)
- `chore`: Maintenance tasks (no version bump)
- `perf`: Performance improvements (triggers patch version bump)
- `ci`: CI/CD changes (no version bump)
- `build`: Build system changes (no version bump)
- `revert`: Revert previous commit (triggers patch version bump)

### Breaking Changes

Add `BREAKING CHANGE:` in the footer or use `!` after the type:

```
feat!: remove deprecated API

BREAKING CHANGE: The old API has been removed
```

## Creating Releases

### Method 1: Automatic (Recommended)

1. Make commits following conventional commit format
2. Push to `main` branch
3. The CI will automatically:
   - Analyze commit messages
   - Determine version bump
   - Create git tag
   - Build and publish release

### Method 2: Manual Trigger

Use GitHub Actions workflow dispatch:

1. Go to **Actions** → **Release** workflow
2. Click **Run workflow**
3. Select version type:
   - `patch`: Bug fixes
   - `minor`: New features
   - `major`: Breaking changes
   - `prerelease`: Pre-release version

### Method 3: Command Line

Use Makefile commands:

```bash
# Patch version (1.0.0 → 1.0.1)
make version-patch

# Minor version (1.0.0 → 1.1.0)
make version-minor

# Major version (1.0.0 → 2.0.0)
make version-major

# Prerelease version
make version-prerelease
```

## Release Process

1. **Version Analysis**: semantic-release analyzes commit messages since last release
2. **Version Bump**: Determines appropriate version bump (major/minor/patch)
3. **Changelog**: Updates `CHANGELOG.md` with new changes
4. **Git Tag**: Creates git tag with new version
5. **Build**: GoReleaser builds binaries for multiple platforms
6. **Publish**: Publishes to GitHub Releases and Docker registry

## Release Artifacts

Each release includes:

- **Binaries**: Linux, Windows, macOS (amd64, arm64)
- **Docker Images**: `ghcr.io/aimar/shelly-prometheus-exporter:latest` and `:vX.Y.Z`
- **Checksums**: SHA256 checksums for all binaries
- **Changelog**: Detailed changelog in release notes

## Examples

### Patch Release

```bash
git commit -m "fix: resolve memory leak in metrics collection"
git push origin main
# Automatically creates v1.0.1
```

### Minor Release

```bash
git commit -m "feat: add support for Shelly 1PM devices"
git push origin main
# Automatically creates v1.1.0
```

### Major Release

```bash
git commit -m "feat!: change configuration format

BREAKING CHANGE: Configuration file format has changed"
git push origin main
# Automatically creates v2.0.0
```

### Prerelease

```bash
git commit -m "feat: add experimental feature"
git push origin main
# Manually trigger prerelease workflow
# Creates v1.1.0-beta.1
```

## Configuration Files

- **`.releaserc.json`**: semantic-release configuration
- **`goreleaser.yml`**: GoReleaser configuration
- **`.github/workflows/release.yml`**: GitHub Actions workflow
- **`package.json`**: Node.js package configuration for semantic-release

## Best Practices

1. **Always use conventional commits** for automatic versioning
2. **Test before releasing** using `make release-snapshot`
3. **Use feature branches** for development
4. **Squash commits** when merging PRs
5. **Write clear commit messages** describing what changed
6. **Use semantic versioning** appropriately for changes

## Troubleshooting

### Release Not Created

- Check commit message format
- Ensure commits are pushed to `main` branch
- Verify GitHub Actions workflow is running

### Wrong Version Bump

- Review commit messages for proper type
- Check if breaking changes are properly marked
- Use manual trigger for specific version type

### Build Failures

- Check GoReleaser configuration
- Verify Docker build process
- Review GitHub Actions logs

## Manual Release (Fallback)

If automatic versioning fails, you can create a manual release:

```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0

# This will trigger the release workflow
```
