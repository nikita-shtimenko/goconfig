package env

// LoaderOption defines a function that configures an Loader
type LoaderOption func(*Loader[any])

// WithSkipMissingFiles allows to not return an error on files that were specified but were not found
func WithSkipMissingFiles() LoaderOption {
	return func(l *Loader[any]) {
		l.skipMissingFiles = true
	}
}
