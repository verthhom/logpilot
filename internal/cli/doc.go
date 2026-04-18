// Package cli provides the command-line interface for logpilot.
//
// It is responsible for:
//   - Parsing flags and positional arguments
//   - Constructing the appropriate Source (file or stdin)
//   - Building the Filter from comma-separated rule strings
//   - Selecting the output format (json or pretty)
//   - Wiring everything together via the pipeline package
//
// Usage:
//
//	logpilot [flags] [file ...]
//
// Flags:
//
//	-format string   output format: json|pretty (default "pretty")
//	-filter string   comma-separated filter rules (e.g. "level=error,msg~timeout")
//
// If no file arguments are provided, logpilot reads from stdin.
package cli
