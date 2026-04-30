package jsongroup

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Grouper batches JSON log lines by the value of a specified field,
// emitting a summary object once the batch window is full.
type Grouper struct {
	field  string
	groups map[string][]json.RawMessage
}

// New creates a Grouper that groups log lines by the given field name.
// Returns an error if field is blank.
func New(field string) (*Grouper, error) {
	if strings.TrimSpace(field) == "" {
		return nil, ErrBlankField
	}
	return &Grouper{
		field:  field,
		groups: make(map[string][]json.RawMessage),
	}, nil
}

// Feed adds a raw JSON line to the appropriate group bucket.
// Non-JSON lines or lines missing the field are silently ignored.
func (g *Grouper) Feed(line string) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return
	}
	raw, ok := obj[g.field]
	if !ok {
		return
	}
	var key string
	if err := json.Unmarshal(raw, &key); err != nil {
		key = string(raw)
	}
	g.groups[key] = append(g.groups[key], json.RawMessage(line))
}

// Flush returns one JSON summary line per group and resets internal state.
// Each summary has the form: {"<field>": "<key>", "count": N, "lines": [...]}
func (g *Grouper) Flush() []string {
	out := make([]string, 0, len(g.groups))
	for key, lines := range g.groups {
		summary := fmt.Sprintf(
			`{"%s":%q,"count":%d,"lines":[%s]}`,
			g.field, key, len(lines), joinRaw(lines),
		)
		out = append(out, summary)
	}
	g.groups = make(map[string][]json.RawMessage)
	return out
}

// Len returns the total number of buffered lines across all groups.
func (g *Grouper) Len() int {
	n := 0
	for _, v := range g.groups {
		n += len(v)
	}
	return n
}

func joinRaw(msgs []json.RawMessage) string {
	parts := make([]string, len(msgs))
	for i, m := range msgs {
		parts[i] = string(m)
	}
	return strings.Join(parts, ",")
}
