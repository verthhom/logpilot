// Package jsoncast provides a Caster that rewrites JSON field values to a
// specified target type (string, int, float, or bool).
//
// Rules are expressed as "field:type" strings, for example:
//
//	"latency:float"
//	"status:int"
//	"retried:bool"
//	"user_id:string"
//
// Non-JSON lines pass through unchanged. Fields that cannot be converted
// are left at their original value.
package jsoncast
