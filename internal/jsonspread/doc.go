// Package jsonspread provides a Spreader that promotes the key-value pairs
// of one or more nested JSON object fields to the top level of the document.
//
// Given a log line such as:
//
//	{"msg":"started","meta":{"host":"srv1","env":"prod"}}
//
// Applying Spreader with field "meta" yields:
//
//	{"msg":"started","host":"srv1","env":"prod"}
//
// Rules:
//   - The source field is deleted after spreading.
//   - If a promoted key already exists at the top level it is left unchanged
//     (first-write-wins semantics).
//   - If the source field is not a JSON object it is left untouched.
//   - Non-JSON lines are passed through without modification.
package jsonspread
