# Agent Guidelines for Airmid Project

This document outlines the guidelines and conventions for agents working on the `airmid` project. Adhering to these guidelines will ensure consistency, maintainability, and efficiency in development.

## General Instructions

* **Understand Context:** Before making any changes, thoroughly understand the existing code, design patterns, and architectural choices within the `airmid` project. Utilize available tools like `read_file`, `search_file_content`, and `list_directory` to gather context.
* **Adhere to Conventions:** Always follow existing project conventions, including naming, formatting, and structural patterns. Mimic the style of surrounding code.
* **Tool Usage:** Prioritize the use of provided tools for tasks like file operations, code modifications, and shell commands.
* **Proactive Testing:** For any feature additions or bug fixes, ensure that appropriate tests are added or updated to cover the changes.
* **Conciseness:** Aim for concise and direct communication. Avoid unnecessary verbosity.

## Go Project Conventions

* **Go Modules:** Dependencies are managed using Go Modules (`go.mod`, `go.sum`). When adding new dependencies, ensure they are properly added and `go mod tidy` is run.
* **File Naming:** Follow standard Go file naming conventions (e.g., `snake_case` for filenames, `_test.go` for test files).
* **Error Handling:** Use Go's idiomatic error handling practices (returning `error` values, checking for `nil` errors).
* **Concurrency:** Use goroutines and channels for concurrency as appropriate, following best practices to avoid race conditions and deadlocks.

## Build and Test

* **Makefile:** The project uses a `Makefile` for common development tasks. Before running custom shell commands, check the `Makefile` for existing targets that accomplish the desired task (e.g., `make test`, `make build`).
* **Testing:**
  * Unit tests are located in `_test.go` files alongside the code they test.
  * Run tests using `make test`.
  * Utilize Go's built-in testing package (`testing`) for writing tests, use `gomega` to assert test values.
* **Linting:** The project uses `golangci-lint` as configured in `.golangci.yaml`. Ensure code adheres to the linting rules. Run `make lint` to check for issues.
* Always run `make test` and `make lint` to ensure correct before commit changes.

## Documentation

* **Code Comments:** Add comments when necessary to explain complex logic, design decisions, or non-obvious behavior. Do not comment on obvious code.
* **API Documentation:** Document public functions, types, and methods using Go's documentation comments.

## Git Practices

* **Commit Messages:** Craft clear, concise, and descriptive commit messages. Focus on *why* a change was made, not just *what* was changed.
* **Branching:** Follow the project's branching strategy (e.g., feature branches, main/master), and you should not working on main branch.

## CI/CD

* **GitHub Workflows:** The `.github/workflows` directory contains GitHub Actions configurations. Be aware of how changes might affect the automated build, test, and deployment processes.
* **Docker:** The `Dockerfile` in `.devcontainer` and the `.dockerignore` file indicate the use of Docker for development environments or deployments.

By following these guidelines, agents can effectively contribute to the `airmid` project while maintaining high code quality and consistency.
