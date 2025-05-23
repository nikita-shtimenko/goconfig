package env

import "github.com/caarlos0/env/v11"

// Options defines a set of functional options for the environment loader
type Options struct {
	SkipMissingFiles bool
	EnvOptions       env.Options
}

// Option defines a functional option for the environment loader
type Option func(*Options) error

// WithSkipMissingFiles configures the loader to skip missing .env files
func WithSkipMissingFiles() Option {
	return func(opts *Options) error {
		opts.SkipMissingFiles = true
		return nil
	}
}

// WithEnvOptions allows passing through options to the underlying env parser
func WithEnvOptions(envOptions env.Options) Option {
	return func(opts *Options) error {
		opts.EnvOptions = envOptions
		return nil
	}
}
