// Package env provides a configuration loader that reads environment variables
// optionally from .env files. It supports generic configuration types and uses
// the github.com/caarlos0/env library for parsing, and github.com/joho/godotenv
// for loading env files.
//
// This package is intended to be used with goconfig to provide environment-based
// configuration loading via a pluggable Loader interface.
package env

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	// ErrEnvFilesNotSpecified indicates that the NewLoader function was called with an empty envFiles array.
	ErrEnvFilesNotSpecified = errors.New("env files not specified")

	// ErrSourceNotFound indicates that the specified source (file, etc.) could not be found.
	ErrSourceNotFound = errors.New("source not found")
)

// Loader implements configuration loading from environment variables
type Loader[T any] struct {
	envFiles         []string
	skipMissingFiles bool
}

// NewLoader creates a new environment-based config loader
func NewLoader[T any](envFiles []string, opts ...LoaderOption) (*Loader[T], error) {
	if len(envFiles) == 0 {
		return nil, ErrEnvFilesNotSpecified
	}

	loader := &Loader[T]{
		envFiles:         envFiles,
		skipMissingFiles: false,
	}

	for _, opt := range opts {
		// This type assertion works because Loader[T] and Loader[any]
		// have the same field layout - we're just changing the type parameter
		typedLoader := (*Loader[any])(unsafe.Pointer(loader))
		opt(typedLoader)
	}

	return loader, nil
}

// Load loads the configuration from environment variables and files
func (l *Loader[T]) Load() (*T, error) {
	// Load environment files using godotenv
	for _, file := range l.envFiles {
		if err := l.loadEnvFile(file); err != nil {
			if l.skipMissingFiles && errors.Is(err, ErrSourceNotFound) {
				continue
			}

			return nil, fmt.Errorf("error loading env file %s: %w", file, err)
		}
	}

	// Parse into struct using caarlos0/env
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		// Just wrap the error with some context - caarlos0/env already provides good error messages
		return nil, fmt.Errorf("error parsing env variables into struct: %w", err)
	}

	return &cfg, nil
}

// loadEnvFile loads environment variables from a .env file using godotenv
func (l *Loader[T]) loadEnvFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return ErrSourceNotFound
	}

	// Use godotenv to load the file
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("failed to load env file: %w", err)
	}

	return nil
}
