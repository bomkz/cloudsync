// Package vtscfg implements encoding and decoding of the block-structured
// config format used by VTOL VR's save files (pilot saves, .vts/.vtm
// scenario/vehicle files, and similar). It is a distinct format from JSON,
// but this package mirrors the encoding/json API as closely as possible so
// it feels familiar: Marshal, Unmarshal, NewEncoder/NewDecoder, struct
// tags, map[string]interface{} decoding, and the Marshaler/Unmarshaler
// escape hatches all work the way their encoding/json counterparts do.
//
// # The format
//
// A file is a sequence of named blocks:
//
//	Name
//	{
//	    key = value
//	    ChildName
//	    {
//	        ...
//	    }
//	}
//
// Compared to JSON, three things are different enough to matter:
//
//  1. There are no quoted strings. A scalar value is simply "everything
//     after the first '=' on the line, trimmed of surrounding whitespace" -
//     so values may freely contain spaces, slashes, semicolons, etc.
//     (e.g. `campaignName = Coop/Freeflight`, `name = a2g gps`).
//  2. A block may legally contain several children with the identical
//     name (e.g. a pilot save has many sibling `VEHICLE { ... }` and
//     `CAMPAIGN { ... }` blocks, and each `currentWeapons` block holds
//     several sibling `weapon { ... }` blocks). This package maps repeated
//     keys onto Go slices; a single occurrence maps onto a scalar/struct
//     field the same way a JSON object field would.
//  3. Vectors/colors are written as a single scalar tuple, e.g.
//     `skinColor = (0.63, 0.48, 0.25)`. Decoding into a []float64, []int,
//     or fixed-size array field auto-detects this. Encoding a slice/array
//     field this way requires the `,tuple` tag option (see below);
//     without it, slice fields are encoded as repeated blocks/lines
//     instead, since that's the more common repetition style in this
//     format.
//
// # Struct tags
//
// Fields are matched using the `vts` struct tag, in the same spirit as
// the `json` tag:
//
//	type Pilot struct {
//	    Name       string    `vts:"pilotName"`
//	    FlightTime float64   `vts:"totalFlightTime"`
//	    SkinColor  []float64 `vts:"skinColor,tuple"`
//	    Vehicles   []Vehicle `vts:"VEHICLE"`
//	}
//
// Supported tag options (comma-separated after the name, name may be left
// blank to keep the Go field name):
//
//   - skip the field entirely
//     omitempty   omit the field when it holds its Go zero value
//     tuple       encode a numeric slice/array as a single "(a, b, c)"
//     scalar instead of repeated blocks/lines
//
// If a field has no `vts` tag, its Go field name is matched first
// exactly and then case-insensitively against block/key names, the same
// fallback encoding/json uses for untagged fields.
//
// # Dynamic / unknown-shape data
//
// Some blocks (notably VDATA) have a shape that depends on which vehicle
// the entry belongs to. Two escape hatches, both modeled on encoding/json:
//
//   - Decode into map[string]interface{} (or an untyped interface{}
//     field) for a generic, JSON-like tree: scalars become string, bool,
//     or float64; a tuple becomes []interface{}; a repeated key becomes
//     []interface{}; anything else becomes a nested map[string]interface{}.
//   - Decode into a Node (or *Node) field to capture a subtree completely
//     unparsed, the way json.RawMessage defers parsing of a JSON value.
//     It can be Marshal'd back out later, or walked manually with
//     Node.Get / Node.GetAll.
//
// # Limitations
//
// This is a best-effort, reflection-based implementation, not a
// byte-for-byte reimplementation of VTOL VR's own (Unity/C#) reader. It
// has no concept of comments or quoting because none appear in real save
// files. Float formatting on Marshal uses Go's shortest round-trippable
// representation, which usually - but not always - matches the exact
// original text byte-for-byte.
package vtscfg

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
