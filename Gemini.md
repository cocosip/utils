# Gemini Project Configuration for go-utils

This document outlines the project-specific guidelines and commands for the `go-utils` library. Adhering to these standards ensures code quality, consistency, and maintainability.

## 1. Project Overview

This is a Go utility library containing various helper functions and modules for common development tasks, including encryption, database interaction, logging, and more.

## 2. Technical Details

- **Go Version**: 1.24
- **Dependency Management**: Go Modules

## 3. Development Workflow

When modifying the codebase, please follow these steps.

### 3.1. Code Formatting

All Go code **must** be formatted using `gofmt`. Before committing any changes, run the following command from the project root, targeting only the files you have changed or added:

```bash
gofmt -w <path-to-your-changed-file-1.go> <path-to-your-changed-file-2.go>
```

### 3.2. Static Analysis

Run `go vet` to catch potential errors and suspicious code structures.

```bash
go vet ./...
```

### 3.3. Dependency Management

After adding or removing dependencies, or after updating the Go version, tidy the `go.mod` and `go.sum` files:

```bash
go mod tidy
```

### 3.4. Running Tests

All new features and bug fixes **must** be accompanied by corresponding unit tests. Ensure all tests pass before committing changes.

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
- `fix(http): Resolve request timeout issue`
- `docs(readme): Update usage instructions for the snowflake package`
- `test(util): Add unit tests for new slice helper functions`
