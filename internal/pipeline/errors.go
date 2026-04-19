package pipeline

import (
	"errors"
	"reflect"
)

// ErrNoSources is returned when a Pipeline is constructed with no sources.
var ErrNoSources = errors.New("pipeline: at least one source is required")

// ErrNilOutput is returned when a Pipeline is constructed with a nil output.
var ErrNilOutput = errors.New("pipeline: output must not be nil")

// Validate checks that the pipeline configuration is valid.
// It returns ErrNoSources if no sources are provided, and ErrNilOutput if out
// is nil or a nil pointer/interface value.
func Validate(sources []Source, out interface{}) error {
	if len(sources) == 0 {
		return ErrNoSources
	}
	if out == nil || (reflect.ValueOf(out).Kind() == reflect.Ptr && reflect.ValueOf(out).IsNil()) {
		return ErrNilOutput
	}
	return nil
}
