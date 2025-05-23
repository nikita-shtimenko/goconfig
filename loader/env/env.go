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

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	// ErrEnvFilesNotSpecified indicates that the NewLoader function was called with an empty Files array.
	ErrEnvFilesNotSpecified = errors.New("env files not specified")

	// ErrSourceNotFound indicates that the specified source (file, etc.) could not be found.
	ErrSourceNotFound = errors.New("source not found")
)

// Loader implements configuration loading from environment variables
type Loader[T any] struct {
	Files   []string
	Options Options
}

// NewLoader creates a new environment-based config loader
func NewLoader[T any](files []string, opts ...Option) (*Loader[T], error) {
	if len(files) == 0 {
		return nil, ErrEnvFilesNotSpecified
	}

	loader := &Loader[T]{
		Files: files,
	}

	for _, opt := range opts {
		if err := opt(&loader.Options); err != nil {
			return nil, fmt.Errorf("error creating loader: invalid option: %w", err)
		}
	}

	return loader, nil
}

// Load loads the configuration from environment variables and files
func (l *Loader[T]) Load() (*T, error) {
	// Load environment files using godotenv
	for _, file := range l.Files {
		if err := l.loadEnvFile(file); err != nil {
			if l.Options.SkipMissingFiles && errors.Is(err, ErrSourceNotFound) {
				continue
			}

			return nil, fmt.Errorf("error loading env file %s: %w", file, err)
		}
	}

	// Parse into struct using caarlos0/env
	var cfg T
	if err := env.ParseWithOptions(&cfg, l.Options.EnvOptions); err != nil {
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
