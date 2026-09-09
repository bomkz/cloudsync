package vtscfg

import (
	"reflect"
	"strings"
)

// isTuple reports whether raw looks like "(a, b, c)".
func isTuple(raw string) bool {
	s := strings.TrimSpace(raw)
	return len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')'
}

// splitTuple splits "(a, b, c)" into []string{"a", "b", "c"}. An empty
// tuple "()" returns nil.
func splitTuple(raw string) []string {
	s := strings.TrimSpace(raw)
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return nil
	}
	parts := strings.Split(inner, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// formatTuple joins parts back into "(a, b, c)".
func formatTuple(parts []string) string {
	return "(" + strings.Join(parts, ", ") + ")"
}

// parseBool parses the format's boolean literals. Game settings files use
// both True/False and Unity-style numeric flags (-1/0/1).
func parseBool(raw string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true":
		return true, true
	case "false":
		return false, true
	case "1", "-1":
		return true, true
	case "0":
		return false, true
	}
	return false, false
}

// formatBool renders a bool the way the format's own files do: "True"/"False".
func formatBool(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
