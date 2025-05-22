# goconfig

[![Go Reference](https://pkg.go.dev/badge/github.com/nikita-shtimenko/goconfig.svg)](https://pkg.go.dev/github.com/nikita-shtimenko/goconfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikita-shtimenko/goconfig)](https://goreportcard.com/report/github.com/nikita-shtimenko/goconfig)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A flexible, generic configuration package for Go applications with a clean strategy pattern design.

## Features

- 🔄 **Strategy Pattern**: Swap configuration sources with zero code changes
- 🧩 **Clean Architecture**: Perfect for dependency injection and testability
- 🧪 **Type-Safe**: Fully leverages Go generics for type safety
- 🛠️ **Extensible**: Easy to create new loaders for different configuration sources
- 🔧 **Options Pattern**: Flexible configuration of loaders using functional options
- 📦 **Minimal Dependencies**: Only imports what's necessary

## Installation

```bash
go get github.com/nikita-shtimenko/goconfig
```

## Quick start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourusername/goconfig"
    "github.com/yourusername/goconfig/loader/env"
)

// Define your configuration structure
type AppConfig struct {
    Server struct {
        Port int `env:"PORT" envDefault:"8080"`
        Host string `env:"HOST" envDefault:"0.0.0.0"`
    }
    Database struct {
        DSN string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/app"`
    }
}

func main() {
    // Create an environment loader
    loader, err := env.NewLoader[AppConfig](
        []string{".env", ".env.local"},
        env.WithSkipMissingFiles(), // Optional: Skip missing files
    )
    if err != nil {
        log.Fatalf("Failed to create config loader: %v", err)
    }
    
    // Load the configuration
    cfg, err := config.NewConfig(loader)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Use the configuration
    fmt.Printf("Server will start on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Database URL: %s\n", cfg.Database.DSN)
}
```

## Usage Guide

### Core concepts

The ```goconfig``` package is built around two core concepts:

1. **ConfigLoader Interface**: A generic interface that any loader must implement
2. **Strategy Pattern**: Different loaders implement the same interface but load from different sources

### Creating a Configuration

1. Define your configuration structure (this example complies with caarl0s/env tags, since we use it for our env loader):

```go
type Config struct {
    LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
    Debug    bool   `env:"DEBUG" envDefault:"false"`
    
    // Nested structures work too
    Server struct {
        Port int `env:"PORT" envDefault:"8080"`
    }
}
```

2. Create a loader for your preferred configuration source:

```go
// Load from environment variables
loader, err := env.NewLoader[Config]([]string{".env"})
```

3. Load the configuration:

```go
cfg, err := config.NewConfig(loader)
```

### Environment Variables Loader

The environment loader can load configuration from:

- Environment variables
- ```env``` files

```go
// Basic usage
loader, err := env.NewLoader[Config]([]string{".env"})

// Multiple .env files (processed in order)
loader, err := env.NewLoader[Config]([]string{".env", ".env.local", ".env.development"})

// With options
loader, err := env.NewLoader[Config](
    []string{".env", ".env.local"},
    env.WithSkipMissingFiles(), // Don't error on missing files
)
```

#### Available Options

- ```WithSkipMissingFiles()```: Skip files that don't exist rather than returning an error

### Extending with Custom Loaders

You can create your own loaders by implementing the ```ConfigLoader[T]``` interface:

```go
// ConfigLoader defines a generic interface for loading configuration
type ConfigLoader[T any] interface {
    Load() (*T, error)
}
```

Example of a custom JSON file loader:

```go
type JsonFileLoader[T any] struct {
    filePath string
}

func NewJsonFileLoader[T any](filePath string) *JsonFileLoader[T] {
    return &JsonFileLoader[T]{filePath: filePath}
}

func (l *JsonFileLoader[T]) Load() (*T, error) {
    data, err := os.ReadFile(l.filePath)
    if err != nil {
        return nil, fmt.Errorf("error reading JSON file: %w", err)
    }
    
    var cfg T
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("error parsing JSON: %w", err)
    }
    
    return &cfg, nil
}
```

## Advanced Usage

### Multiple Configuration Sources

You can implement custom loaders that combine multiple sources, or load configurations separately and combine them in your application:

```go
// Load service-specific configuration
serviceCfg, err := config.NewConfig(serviceLoader)

// Load database configuration
dbCfg, err := config.NewConfig(dbLoader)

// Use them together
app := NewApp(serviceCfg, dbCfg)
```

### Composing Configurations from Different Sources

A recommended pattern is to create domain-specific configurations and compose them together:

```go
// Define domain-specific config structures
type FooConfig struct {
    Foo string `env:"FOO,required"`
}

type BarConfig struct {
    Bar string `env:"BAR,required"`
}

// Create a composite config structure
type Config struct {
    Foo *FooConfig
    Bar *BarConfig
}

func main() {
    fooConfigLoader, err := env.NewLoader[FooConfig](
        []string{".env", ".env.local"},
        env.WithSkipMissingFiles(),
    )

    if err != nil {
        log.Fatal(err)
    }

    fooConfig, err := config.NewConfig(fooConfigLoader)
    if err != nil {
        log.Fatal(err)
    }

    barConfigLoader, err := env.NewLoader[BarConfig](
        []string{".env"},
    )

    if err != nil {
        log.Fatal(err)
    }

    barConfig, err := config.NewConfig(barConfigLoader)
    if err != nil {
        log.Fatal(err)
    }

    // Compose the configurations
    cfg := &Config{
        Foo: fooConfig,
        Bar: barConfig,
    }

    // Use the composite configuration
    // ...
}

```

This pattern offers several benefits:

- Clear separation of configuration domains
- Type-safety for each component
- Ability to use different sources for different parts of your config
- Explicit composition of the final configuration
- Better testability of individual components

## Built-in loaders

1. **env** - environment loader (loads from .env files)

## License

MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
