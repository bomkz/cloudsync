package vtscfg

// Node represents one parsed block:
//
//	Name
//	{
//	    ...Fields...
//	}
//
// The root Node returned by Parse is synthetic (Name == "") and its
// Fields are the file's top-level blocks (normally just one, e.g.
// "PILOTS").
//
// Node is the low-level, order-preserving, duplicate-key-preserving
// analogue of the map/slice tree encoding/json builds internally. Most
// callers won't need to touch it directly - Marshal/Unmarshal build and
// consume it automatically - but it's exported for advanced use (walking
// a dynamic subtree by hand) and so a struct field of type Node/*Node can
// capture a subtree unparsed, the way json.RawMessage does for JSON.
type Node struct {
	Name   string
	Fields []Field
}

// Field is one entry inside a Node. Exactly one of Value or Child is set:
// a scalar assignment ("key = value") sets Value; a nested block
// ("key { ... }") sets Child.
type Field struct {
	Key   string
	Value *string
	Child *Node
}

// IsScalar reports whether this field is a "key = value" assignment.
func (f Field) IsScalar() bool { return f.Value != nil }

// IsNode reports whether this field is a nested "key { ... }" block.
func (f Field) IsNode() bool { return f.Child != nil }

// RawValue returns the field's scalar text, or "" if it is a nested node.
func (f Field) RawValue() string {
	if f.Value == nil {
		return ""
	}
	return *f.Value
}

// Get returns the first field with the given key.
func (n *Node) Get(key string) (Field, bool) {
	if n == nil {
		return Field{}, false
	}
	for _, f := range n.Fields {
		if f.Key == key {
			return f, true
		}
	}
	return Field{}, false
}

// GetAll returns every field with the given key, in file order.
func (n *Node) GetAll(key string) []Field {
	if n == nil {
		return nil
	}
	var out []Field
	for _, f := range n.Fields {
		if f.Key == key {
			out = append(out, f)
		}
	}
	return out
}

// Value returns the raw scalar text for the first field with the given
// key, or "" if there is no such field (or it is a nested node).
func (n *Node) Value(key string) string {
	f, ok := n.Get(key)
	if !ok {
		return ""
	}
	return f.RawValue()
}

// Child returns the first nested node with the given key.
func (n *Node) Child(key string) (*Node, bool) {
	f, ok := n.Get(key)
	if !ok || f.Child == nil {
		return nil, false
	}
	return f.Child, true
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
