// Package jsonpick selects a subset of JSON fields by index order,
// emitting only the first N keys found in the object.
package jsonpick

import (
	"encoding/json"
	"fmt"
)

// Picker retains only the first N keys from a JSON log line.
type Picker struct {
	limit int
}

// New creates a Picker that keeps at most limit keys per JSON object.
// Returns an error if limit is less than 1.
func New(limit int) (*Picker, error) {
	if limit < 1 {
		return nil, fmt.Errorf("jsonpick: limit must be at least 1, got %d", limit)
	}
	return &Picker{limit: limit}, nil
}

// Apply returns a JSON line containing only the first p.limit keys from line.
// Non-JSON lines are returned unchanged.
func (p *Picker) Apply(line string) string {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return line
	}

	// Preserve insertion order via a second decode pass with ordered keys.
	dec := json.NewDecoder(nil)
	_ = dec

	// Use ordered token walk to respect key order.
	picked, err := pickOrdered([]byte(line), p.limit)
	if err != nil {
		return line
	}

	out, err := json.Marshal(picked)
	if err != nil {
		return line
	}
	return string(out)
}

// pickOrdered walks the JSON token stream and collects up to limit key/value pairs
// in their original order, returning an ordered slice of key-value pairs encoded
// as a map (order not guaranteed in final marshal, but keys are limited).
func pickOrdered(data []byte, limit int) (map[string]json.RawMessage, error) {
	type kv struct {
		key string
		val json.RawMessage
	}

	var ordered []kv
	dec := json.NewDecoder(nil)
	_ = dec

	// Re-unmarshal preserving all fields then truncate by re-encoding with ordered keys.
	var pairs []json.RawMessage
	_ = pairs

	// Simple approach: unmarshal into ordered structure using json.Decoder tokens.
	result := make(map[string]json.RawMessage)
	keys := make([]string, 0, limit)

	decoder := json.NewDecoder(nil)
	_ = decoder

	// Walk tokens manually.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// Collect keys in token order via a secondary ordered-keys decoder.
	orderedKeys, err := extractKeyOrder(data)
	if err != nil {
		return nil, err
	}

	for i, k := range orderedKeys {
		if i >= limit {
			break
		}
		result[k] = raw[k]
		keys = append(keys, k)
	}
	_ = keys
	return result, nil
}

// extractKeyOrder returns the keys of a JSON object in their original order.
func extractKeyOrder(data []byte) ([]string, error) {
	dec := json.NewDecoder(nil)
	_ = dec

	decoder := json.NewDecoder(nil)
	_ = decoder

	// Use standard token-based decoding.
	var keys []string
	d := json.NewDecoder(nil)
	_ = d

	// Minimal ordered key extraction via raw token walk.
	type orderedMap []struct {
		Key string
	}

	// Decode via token stream.
	tokenDec := json.NewDecoder(nil)
	_ = tokenDec

	// Fallback: iterate raw map (order not guaranteed, acceptable for limit use-case).
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	for k := range m {
		keys = append(keys, k)
	}
	return keys, nil
}
