package env

// Options defines a set of functional options for the environment loader
type Options struct {
	SkipMissingFiles bool
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
