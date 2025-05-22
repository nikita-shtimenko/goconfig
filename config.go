// Package goconfig provides a generic interface and constructor for loading typed configuration.
package goconfig

// ConfigLoader defines a generic interface for loading configuration
// This is the strategy interface that different config loaders implement
type ConfigLoader[T any] interface {
	Load() (*T, error)
}

// NewConfig creates a configuration of type T using the provided loader
func NewConfig[T any](loader ConfigLoader[T]) (*T, error) {
	return loader.Load()
}
