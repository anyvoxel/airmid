<h1 align="center">Airmid</h1>

<p align="center">
  <strong>A lightweight, extensible, and production-ready Inversion of Control (IoC) framework for Go.</strong>
</p>

<p align="center">
    <a href="https://github.com/anyvoxel/airmid/actions/workflows/tests.yml">
        <img src="https://github.com/anyvoxel/airmid/actions/workflows/tests.yml/badge.svg" alt="Build Status">
    </a>
    <a href="https://codecov.io/gh/anyvoxel/airmid">
        <img src="https://codecov.io/gh/anyvoxel/airmid/branch/main/graph/badge.svg" alt="Code Coverage">
    </a>
    <a href="https://goreportcard.com/report/github.com/anyvoxel/airmid">
        <img src="https://goreportcard.com/badge/github.com/anyvoxel/airmid" alt="Go Report Card">
    </a>
    <a href="https://github.com/anyvoxel/airmid/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/anyvoxel/airmid" alt="License">
    </a>
    <a href="https://pkg.go.dev/github.com/anyvoxel/airmid">
        <img src="https://pkg.go.dev/badge/github.com/anyvoxel/airmid.svg" alt="Go Reference">
    </a>
</p>

Airmid is a powerful Inversion of Control (IoC) container for Go, designed to simplify dependency management and application composition. Inspired by the principles of the Spring Framework, Airmid helps you build modular, testable, and maintainable applications by managing the lifecycle and wiring of your application components (beans).

## Features

- **Dependency Injection:** Automatically inject dependencies into your structs, simplifying initialization and decoupling components.
- **Bean Lifecycle Management:** Full support for singleton and prototype scopes, with hooks for initialization (`InitializingBean`) and destruction (`DestructionAwareBeanPostProcessor`).
- **Extensible Architecture:** Customize the framework's behavior with `BeanPostProcessor` and `BeanDefinitionPostProcessor`.
- **Configuration Management:** Externalize configuration from your code using property files, with support for profiles and environment variables.
- **Declarative Component Registration:** Define beans and their dependencies in a declarative way, cleaning up your `main` function and application startup logic.

## Getting Started

### Installation

To start using Airmid, install the necessary packages in your project:

```sh
go get github.com/anyvoxel/airmid/app
go get github.com/anyvoxel/airmid/ioc
```

### Quick Example

Here’s a simple example to demonstrate the basic functionality of Airmid.

1. **Define a component (Bean) and Runner:**

    Create a struct that will act as a service. We'll use an `init()` function to register it with the Airmid container.

    `main.go`:

    ```go
    package main

    type selfString string

    type runner3 struct {
        i selfString `airmid:"value:${Airmid}"`
    }

    func (r *runner3) Run(ctx context.Context) {
        fmt.Println("Hello, ", r.i, "!")
    }

    func (r *runner3) Stop(ctx context.Context) {

    }

    func init() {
        anvil.Must(
            app.RegisterBeanDefinition("runner3", ioc.MustNewBeanDefinition(reflect.TypeOf((*runner3)(nil)))),
        )
    }
    ```

1. **Run the application:**

    The `main` function is now incredibly simple. It just needs to start the Airmid application.

    `main.go`:

    ```go
    package main

    import (
        "context"
        "log"
        
        "github.com/anyvoxel/airmid/app"
    )
    
    func main() {
        // Launch the application. Airmid finds and runs all ApplicationRunner beans.
        if err := app.Run(context.Background()); err != nil {
            log.Fatalf("Failed to run application: %v", err)
        }
    }
    ```

1. **Run it!**

    When you run `go run .`, you will see the output:

    ```
    Hello, Airmid!
    ```

## Building From Source

To contribute to Airmid, you can build it from the source.

1. **Clone the repository:**

    ```sh
    git clone https://github.com/anyvoxel/airmid.git
    cd airmid
    ```

2. **Run tests:**
    To ensure everything is working correctly, run the full test suite.

    ```sh
    make test
    ```

3. **Run linter:**
    We use `golangci-lint` to maintain code quality.

    ```sh
    make lint
    ```

## Contributing

We welcome contributions from the community! Please read our `CONTRIBUTING.md` (coming soon) for guidelines on how to contribute. You can start by checking our issue tracker for open tasks.

## License

Airmid is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
