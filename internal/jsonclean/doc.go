// Package jsonclean provides a Cleaner that strips unwanted zero-value fields
// from structured JSON log lines.
//
// Supported removal targets are controlled via Options:
//
//	- null values
//	- empty strings ("")
//	- empty arrays ([])
//	- empty objects ({})
//
// Non-JSON lines are passed through unchanged so that the cleaner can be
// safely inserted into any pipeline stage without breaking plain-text sources.
//
// Example:
//
//	c := jsonclean.New(jsonclean.Options{
//		RemoveNull:        true,
//		RemoveEmptyString: true,
//	})
//	out := c.Apply(`{"level":"info","err":null,"msg":""}`) // → {"level":"info"}
package jsonclean
