# Qwen Project Configuration for go-utils

This document outlines the project-specific guidelines and commands for the `go-utils` library. Adhering to these standards ensures code quality, consistency, and maintainability.

## 1. Project Overview

This is a Go utilities library containing various helper functions and modules for common development tasks, including cryptography, database interactions, logging, and more.

## 2. Technical Details

- **Go Version**: 1.24
- **Dependency Management**: Go Modules

## 3. Development Workflow

Follow these steps when making changes to the codebase.

### 3.1. Code Formatting

All Go code **MUST** be formatted using `gofmt`. Before committing any changes, run the following command from the project root to format all files:

```bash
gofmt -w .
```

### 3.2. Static Analysis

Run `go vet` to catch potential errors and suspicious code constructs.

```bash
go vet ./...
```

### 3.3. Dependency Management

After adding or removing a dependency, or updating the Go version, tidy the `go.mod` and `go.sum` files:

```bash
go mod tidy
```

### 3.4. Running Tests

All new features and bug fixes **MUST** be accompanied by corresponding unit tests. Ensure all tests pass before submitting changes.

To run all tests in the project, use the following command:

```bash
go test -v ./...
```

## 4. Committing Changes

Commit messages should be clear and follow a consistent format. A good practice is to reference the package being changed.

**Format:**
`<type>(<scope>): <subject>`

**Examples:**

- `feat(crypto): Add support for AES-GCM encryption`
- `fix(http): Resolve issue with request timeouts`
- `docs(readme): Update usage instructions for the snowflake package`
- `test(util): Add unit tests for new slice helper functions`