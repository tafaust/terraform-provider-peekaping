# Contributing to Terraform Provider for Peekaping

Thank you for your interest in contributing to the Terraform Provider for Peekaping! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style and Standards](#code-style-and-standards)
- [Documentation](#documentation)
- [Commit Messages](#commit-messages)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Release Process](#release-process)
- [Reporting Issues](#reporting-issues)
- [Community Support](#community-support)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating in this project, you agree to abide by its terms.

## Getting Started

### Prerequisites

Before contributing, ensure you have the following installed:

- **Go 1.24+**: [Install Go](https://golang.org/doc/install)
- **Terraform 1.7+**: [Install Terraform](https://www.terraform.io/downloads)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **asdf** (recommended): [Install asdf](https://asdf-vm.com/guide/getting-started.html)

### Fork and Clone

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/your-username/terraform-provider-peekaping.git
   cd terraform-provider-peekaping
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/tafaust/terraform-provider-peekaping.git
   ```

## Development Environment

### Using asdf (Recommended)

This project uses asdf for tool version management. The required versions are specified in `.tool-versions`:

```bash
# Install asdf if not already installed
# Then install the required tools
asdf install
```

### Manual Setup

If you prefer not to use asdf:

1. **Install Go 1.24.2** or later
2. **Install Terraform 1.7.5** or later
3. **Install development tools**:
   ```bash
   cd tools
   go generate
   ```

### Environment Variables

For testing and development, you may need to set these environment variables:

```bash
export PEEKAPING_ENDPOINT="https://api.peekaping.com"  # or your instance URL
export PEEKAPING_EMAIL="your-email@example.com"
export PEEKAPING_PASSWORD="your-password"
export PEEKAPING_TOKEN="123456"  # 2FA token if enabled
```

## Project Structure

```
├── internal/
│   ├── peekaping/          # API client implementation
│   └── provider/           # Terraform provider implementation
│       ├── datasource_*.go # Data source implementations
│       └── resource_*.go   # Resource implementations
├── examples/               # Example configurations
│   ├── simple/            # Basic usage examples
│   ├── comprehensive/     # Advanced usage examples
│   └── full-lifecycle/    # Complete lifecycle examples
├── docs/                  # Generated documentation
├── tools/                 # Development tools
├── .github/workflows/     # CI/CD workflows
├── GNUmakefile           # Build and development commands
└── main.go               # Provider entry point
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### 2. Make Your Changes

- Follow the [Code Style and Standards](#code-style-and-standards)
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests

```bash
# Run unit tests
make test

# Run tests with race detection
make testrace

# Run tests with coverage
make testcover

# Run acceptance tests (requires valid credentials)
make testacc
```

### 4. Format and Lint

```bash
# Format code
make fmt

# Run linters
make lint

# Run all checks
make check
```

### 5. Commit Your Changes

Follow the [Commit Messages](#commit-messages) guidelines.

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Testing

### Unit Tests

Unit tests are located alongside the code they test. Run them with:

```bash
make test
```

### Acceptance Tests

Acceptance tests create real resources and require valid credentials:

```bash
# Set required environment variables
export PEEKAPING_ENDPOINT="https://api.peekaping.com"
export PEEKAPING_EMAIL="your-email@example.com"
export PEEKAPING_PASSWORD="your-password"
export PEEKAPING_TOKEN="123456"

# Run acceptance tests
make testacc
```

### Test Coverage

Generate and view test coverage:

```bash
make testcover
# Opens coverage.html in your browser
```

### Example Validation

Validate all example configurations:

```bash
make validate-examples
```

## Code Style and Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (enforced by CI)
- Follow the project's linting rules (see `.golangci.yml`)

### Terraform Code

- Use consistent formatting with `terraform fmt`
- Follow [Terraform style conventions](https://www.terraform.io/docs/language/syntax/style.html)

### Naming Conventions

- **Resources**: `peekaping_<resource_type>` (e.g., `peekaping_monitor`)
- **Data Sources**: `peekaping_<resource_type>` (e.g., `peekaping_monitor`)
- **Provider**: `peekaping`

### Code Organization

- Keep API client code in `internal/peekaping/`
- Keep provider implementation in `internal/provider/`
- Use descriptive function and variable names
- Add comments for exported functions and types
- Follow the existing patterns in the codebase

## Documentation

### Code Documentation

- Add Go doc comments for all exported functions and types
- Use clear, concise descriptions
- Include examples where helpful

### Terraform Documentation

Documentation is generated automatically using `tfplugindocs`:

```bash
make docs
```

This generates documentation in the `docs/` directory based on schema definitions.

### Example Documentation

- Keep examples in the `examples/` directory
- Add README files for complex examples
- Ensure examples are working and tested

## Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(monitor): add support for gRPC monitors
fix(provider): handle 2FA token expiration gracefully
docs: update installation instructions
test(monitor): add acceptance tests for HTTP monitors
```

## Pull Request Guidelines

### Before Submitting

1. **Run all checks**:
   ```bash
   make check
   ```

2. **Ensure tests pass**:
   ```bash
   make test
   ```

3. **Update documentation** if needed:
   ```bash
   make docs
   ```

4. **Rebase on main**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### PR Description

Include the following in your PR description:

- **Summary**: Brief description of changes
- **Type**: Bug fix, feature, documentation, etc.
- **Testing**: How you tested the changes
- **Breaking Changes**: Any breaking changes (if applicable)
- **Related Issues**: Link to related issues

### PR Requirements

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
- [ ] Commit messages follow conventional format

## Release Process

### Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Steps

1. **Update version** in `main.go`
2. **Update CHANGELOG.md** with new features/fixes
3. **Create release tag**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
4. **Create GitHub release** with release notes

## Reporting Issues

### Bug Reports

When reporting bugs, include:

- **Terraform version**: `terraform version`
- **Provider version**: Version of the provider
- **Go version**: `go version` (if building from source)
- **Steps to reproduce**: Clear, minimal steps
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Error messages**: Full error output
- **Configuration**: Relevant Terraform configuration (sanitized)

### Feature Requests

For feature requests, include:

- **Use case**: Why is this feature needed?
- **Proposed solution**: How should it work?
- **Alternatives**: Other solutions considered
- **Additional context**: Any other relevant information

## Community Support

### Getting Help

- **Documentation**: Check the [docs/](docs/) directory
- **Examples**: Look at [examples/](examples/) for usage patterns
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions

### Contributing Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Provide constructive feedback
- Follow the code of conduct

## Development Tools

### Available Make Targets

```bash
make help  # Show all available targets
```

Key targets:
- `make build` - Build the provider
- `make test` - Run unit tests
- `make testacc` - Run acceptance tests
- `make lint` - Run linters
- `make fmt` - Format code
- `make docs` - Generate documentation
- `make check` - Run all checks
- `make clean` - Clean build artifacts

### IDE Setup

Recommended VS Code extensions:
- Go extension
- Terraform extension
- GitLens

### Debugging

For debugging Terraform provider issues:

1. Enable debug logging:
   ```bash
   export TF_LOG=DEBUG
   export TF_LOG_PATH=terraform.log
   ```

2. Use the provider in debug mode:
   ```bash
   terraform init
   terraform plan
   ```

## License

By contributing to this project, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

Thank you for contributing to the Terraform Provider for Peekaping! 🎉
