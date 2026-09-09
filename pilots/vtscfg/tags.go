package vtscfg

import "strings"

// tagOptions is the parsed form of a `vts:"..."` struct tag.
type tagOptions struct {
	name      string
	omitempty bool
	tuple     bool
	skip      bool
}

// parseTag parses a `vts` struct tag value (as returned by
// reflect.StructTag.Lookup) for a field whose Go name is fieldName. An
// empty tag (no `vts` tag present) keeps the Go field name and no
// options. A tag of "-" means skip.
func parseTag(raw string, hasTag bool, fieldName string) tagOptions {
	if hasTag && raw == "-" {
		return tagOptions{skip: true}
	}
	opt := tagOptions{name: fieldName}
	if raw == "" {
		return opt
	}
	parts := strings.Split(raw, ",")
	if parts[0] != "" {
		opt.name = parts[0]
	}
	for _, p := range parts[1:] {
		switch strings.TrimSpace(p) {
		case "omitempty":
			opt.omitempty = true
		case "tuple":
			opt.tuple = true
		}
	}
	return opt
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
