# CodeClimate Setup Guide

This guide explains how to set up CodeClimate for the Shelly Prometheus Exporter project.

## What is CodeClimate?

CodeClimate is a code quality analysis tool that provides:

- **Code quality metrics**: Maintainability, complexity, and style analysis
- **Test coverage tracking**: Monitor test coverage over time
- **Security analysis**: Identify potential security vulnerabilities
- **Duplication detection**: Find duplicated code patterns
- **Technical debt tracking**: Monitor code quality trends

## Setup Steps

### 1. Create CodeClimate Account

1. Go to [CodeClimate](https://codeclimate.com/)
2. Sign up with your GitHub account
3. Connect your GitHub repositories

### 2. Add Repository to CodeClimate

1. In CodeClimate dashboard, click "Add a repository"
2. Select `aimarjs/shelly-prometheus-exporter`
3. Choose "Open Source" plan (free for public repositories)
4. Wait for the initial analysis to complete

### 3. Get Coverage Token

**Note**: For coverage reporting, you now need a Qlty coverage token instead of the old CodeClimate Test Reporter ID.

#### Option 1: Qlty Coverage Token

1. Go to your Qlty dashboard (https://qlty.sh)
2. Navigate to your project settings
3. Look for "Coverage" or "Code Coverage" section
4. Copy the coverage token (e.g., `qltcp_eho0pwW0HM4RtkJr`)
5. Add it as a GitHub secret: `QLTY_COVERAGE_TOKEN`

#### Option 2: Skip Test Coverage (Alternative)

If you don't want to use Qlty for coverage, you can skip the coverage reporting and just use CodeClimate for code quality analysis.

### 4. Configure GitHub Secret

1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `QLTY_COVERAGE_TOKEN`
5. Value: Your Qlty coverage token

## Configuration Files

### `.codeclimate.yml`

The main configuration file that defines:

- Which files to exclude from analysis
- Which engines to enable
- Quality thresholds and checks
- Plugin configurations

### Key Configuration Options

#### Excluded Files

```yaml
exclude_patterns:
  - "**/*_test.go" # Test files
  - "**/vendor/**" # Dependencies
  - "**/docs/**" # Documentation
  - "**/.github/**" # GitHub workflows
```

#### Enabled Engines

```yaml
engines:
  gofmt: # Code formatting
  golint: # Linting
  govet: # Go vet analysis
  gosec: # Security analysis
  go-cyclo: # Cyclomatic complexity
  duplication: # Code duplication
  maintainability: # Maintainability rating
```

#### Quality Thresholds

```yaml
checks:
  function-length:
    config:
      threshold: 50 # Max function length
  file-length:
    config:
      threshold: 300 # Max file length
  argument-count:
    config:
      threshold: 4 # Max function arguments
```

## GitHub Actions Integration

### Workflow Configuration

The project includes both CodeClimate analysis and Qlty coverage reporting in the CI pipeline:

```yaml
codeclimate:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Set up Go
      uses: actions/setup-go@v4
    - name: CodeClimate Analysis
      uses: paambaati/codeclimate-action@v3.2.0
      with:
        coverageCommand: make test-coverage
        coverageLocations: coverage.out:lcov
      env:
        CC_TEST_REPORTER_ID: ${{ secrets.CC_TEST_REPORTER_ID }}

test:
  runs-on: ubuntu-latest
  steps:
    # ... other steps ...

    - name: Run coverage
      run: CGO_ENABLED=0 go test -coverprofile=coverage.out -covermode=atomic ./...

    - name: Upload coverage to Qlty
      uses: qltysh/qlty-action/coverage@v2
      with:
        token: ${{ secrets.QLTY_COVERAGE_TOKEN }}
        files: coverage.out
```

### Coverage Reporting

The workflow reports test coverage to both CodeClimate and Qlty:

- Runs tests with coverage
- Generates LCOV format coverage report
- Uploads coverage data to CodeClimate (for code quality analysis)
- Uploads coverage data to Qlty (for coverage reporting)

## Badges and Status

### README Badges

The following badges are added to the README:

```markdown
[![Code Coverage](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter/coverage.svg)](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter)
[![Maintainability](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter/maintainability.svg)](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter)
```

### Badge Types

- **GPA Badge**: Overall code quality grade (A-F)
- **Coverage Badge**: Test coverage percentage
- **Maintainability Badge**: Maintainability rating

## Quality Metrics

### Maintainability Rating

CodeClimate provides a maintainability rating based on:

- **Cyclomatic Complexity**: Code complexity analysis
- **Duplication**: Code duplication detection
- **Style Issues**: Code style violations
- **Security Issues**: Potential security vulnerabilities

### Test Coverage

Tracks test coverage over time:

- **Line Coverage**: Percentage of code lines tested
- **Branch Coverage**: Percentage of code branches tested
- **Function Coverage**: Percentage of functions tested

### Security Analysis

Identifies potential security issues:

- **gosec**: Go security analysis
- **Common vulnerabilities**: Known security patterns
- **Dependency issues**: Vulnerable dependencies

## Best Practices

### Code Quality

1. **Keep functions small**: Aim for < 50 lines
2. **Limit complexity**: Keep cyclomatic complexity < 10
3. **Avoid duplication**: Refactor common patterns
4. **Write tests**: Maintain high test coverage

### Configuration

1. **Exclude non-source files**: Don't analyze docs, configs, etc.
2. **Set appropriate thresholds**: Balance quality vs. practicality
3. **Monitor trends**: Watch quality metrics over time
4. **Address issues promptly**: Fix quality issues as they arise

### Workflow Integration

1. **Run on every PR**: Ensure quality checks on all changes
2. **Fail on quality regression**: Block PRs that reduce quality
3. **Monitor coverage**: Maintain or improve test coverage
4. **Review reports**: Regularly review CodeClimate reports

## Troubleshooting

### Common Issues

#### Analysis Not Running

- Check GitHub secret `CC_TEST_REPORTER_ID` is set
- Verify CodeClimate repository is connected
- Ensure workflow has proper permissions

#### Coverage Not Reporting

- Verify `coverage.out` file is generated
- Check LCOV format is correct
- Ensure test coverage command runs successfully

#### Quality Issues

- Review CodeClimate report for specific issues
- Check threshold configurations
- Consider excluding false positives

### Debug Steps

1. **Check GitHub Actions logs**: Look for CodeClimate step output
2. **Verify secrets**: Ensure `CC_TEST_REPORTER_ID` is set
3. **Test locally**: Run coverage command manually
4. **Review configuration**: Check `.codeclimate.yml` syntax

## Advanced Configuration

### Custom Engines

You can add custom analysis engines:

```yaml
engines:
  custom-engine:
    enabled: true
    config:
      threshold: 5
```

### Plugin Configuration

Configure specific plugins:

```yaml
plugins:
  go-cyclo:
    enabled: true
    config:
      threshold: 10
```

### Quality Gates

Set up quality gates to fail builds:

```yaml
checks:
  maintainability:
    enabled: true
    config:
      threshold: 4.0
```

## Monitoring and Alerts

### Quality Trends

Monitor code quality over time:

- **Maintainability trends**: Track quality improvements
- **Coverage trends**: Monitor test coverage changes
- **Issue resolution**: Track issue fix rates

### Alerts

Set up alerts for:

- **Quality regression**: When quality drops
- **Coverage decrease**: When test coverage falls
- **New issues**: When new quality issues appear

## Resources

- [CodeClimate Documentation](https://docs.codeclimate.com/)
- [Go Analysis Engines](https://docs.codeclimate.com/docs/go)
- [Test Coverage Setup](https://docs.codeclimate.com/docs/test-coverage)
- [GitHub Actions Integration](https://docs.codeclimate.com/docs/github-actions)
