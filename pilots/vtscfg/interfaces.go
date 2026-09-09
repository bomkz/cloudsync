package vtscfg

// Unmarshaler is implemented by types that want full control over how a
// nested block is decoded, the same role json.Unmarshaler plays for
// encoding/json. UnmarshalVTS receives the already-parsed subtree; it
// must copy any data it needs, since the Node may be reused/mutated.
type Unmarshaler interface {
	UnmarshalVTS(*Node) error
}

// Marshaler is implemented by types that want full control over how they
// are encoded as a nested block, the same role json.Marshaler plays for
// encoding/json. The returned Node's Name is overwritten by the caller
// with the field/map/slice key it's being encoded under, so it need not
// be set.
type Marshaler interface {
	MarshalVTS() (*Node, error)
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
