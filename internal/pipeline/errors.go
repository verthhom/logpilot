package pipeline

import "errors"

// ErrNoSources is returned when a Pipeline is constructed with no sources.
var ErrNoSources = errors.New("pipeline: at least one source is required")

// ErrNilOutput is returned when a Pipeline is constructed with a nil output.
var ErrNilOutput = errors.New("pipeline: output must not be nil")

// Validate checks that the pipeline configuration is valid.
func Validate(sources []Source, out interface{}) error {
	if len(sources) == 0 {
		return ErrNoSources
	}
	if out == nil {
		return ErrNilOutput
	}
	return nil
}
