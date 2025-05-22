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

// ErrEnvFilesNotSpecified indicates that call to NewLoader function were made with empty envFiles array
var ErrEnvFilesNotSpecified = errors.New("env files not specified")

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

// Load loads environment variables into the generic configuration type
func (l *Loader[T]) Load() (*T, error) {
	var cfg T

	for _, file := range l.envFiles {
		// Check if file exists before attempting to load
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if l.skipMissingFiles {
				continue
			}

			return nil, fmt.Errorf("error loading env file %s: %w", file, err)
		}

		if err := godotenv.Load(file); err != nil {
			return nil, fmt.Errorf("error loading env file %s: %w", file, err)
		}
	}

	// Parse environment variables into the config struct
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing env vars: %w", err)
	}

	return &cfg, nil
}
