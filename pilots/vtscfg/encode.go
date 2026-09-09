package vtscfg

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
)

// Marshal returns the config-format encoding of v, the same way
// json.Marshal does for JSON. v (after dereferencing pointers/interfaces)
// must be a struct or map[string]X; its top-level fields become the
// file's top-level blocks. Output is indented with tabs, matching the
// style VTOL VR's own files use.
func Marshal(v interface{}) ([]byte, error) {
	return MarshalIndent(v, "", "\t")
}

// MarshalIndent is like Marshal but applies indent at each nesting level,
// prefixed by prefix, mirroring json.MarshalIndent.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	fields, err := encodeTopLevel(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	writeFields(&buf, fields, prefix, indent)
	return buf.Bytes(), nil
}

// Encoder writes config-format documents to an output stream, mirroring
// json.Encoder's constructor/method shape.
type Encoder struct {
	w      io.Writer
	prefix string
	indent string
}

// NewEncoder returns a new Encoder that writes to w. The default
// indentation is a single tab per nesting level, matching VTOL VR's own
// files; use SetIndent to change it.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, indent: "\t"}
}

// SetIndent configures the prefix and per-level indent used by Encode,
// mirroring json.Encoder.SetIndent.
func (e *Encoder) SetIndent(prefix, indent string) {
	e.prefix = prefix
	e.indent = indent
}

// Encode writes the config-format encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	b, err := MarshalIndent(v, e.prefix, e.indent)
	if err != nil {
		return err
	}
	_, err = e.w.Write(b)
	return err
}

func encodeTopLevel(v interface{}) ([]Field, error) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, nil
	}
	rv = derefForEncode(rv)
	if !rv.IsValid() {
		return nil, nil
	}
	switch rv.Kind() {
	case reflect.Struct:
		if rv.Type() == nodeType {
			n := rv.Interface().(Node)
			return n.Fields, nil
		}
		return encodeStructFields(rv)
	case reflect.Map:
		return encodeMapFields(rv)
	default:
		return nil, &UnsupportedTypeError{Type: rv.Type()}
	}
}

// derefForEncode walks a chain of pointers, returning an invalid Value
// (signaling "omit entirely") if it hits a nil pointer.
func derefForEncode(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

func marshalerFor(v reflect.Value) (Marshaler, bool) {
	if !v.IsValid() {
		return nil, false
	}
	if m, ok := v.Interface().(Marshaler); ok {
		return m, true
	}
	if v.CanAddr() {
		if m, ok := v.Addr().Interface().(Marshaler); ok {
			return m, true
		}
	}
	return nil, false
}

func textMarshalerFor(v reflect.Value) (encoding.TextMarshaler, bool) {
	if !v.IsValid() {
		return nil, false
	}
	if tm, ok := v.Interface().(encoding.TextMarshaler); ok {
		return tm, true
	}
	if v.CanAddr() {
		if tm, ok := v.Addr().Interface().(encoding.TextMarshaler); ok {
			return tm, true
		}
	}
	return nil, false
}

// appendEncodedField encodes v (one Go value) under key, appending
// one-or-more Fields onto *out. tuple requests "(a, b, c)" scalar
// encoding for numeric slice/array values instead of the default
// repeated-entries encoding.
func appendEncodedField(out *[]Field, key string, v reflect.Value, tuple bool) error {
	v = derefForEncode(v)
	if !v.IsValid() {
		return nil // nil pointer: omitted, since there's no null literal
	}

	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		return appendEncodedField(out, key, v.Elem(), tuple)
	}

	if m, ok := marshalerFor(v); ok {
		n, err := m.MarshalVTS()
		if err != nil {
			return err
		}
		if n == nil {
			return nil
		}
		nCopy := *n
		nCopy.Name = key
		*out = append(*out, Field{Key: key, Child: &nCopy})
		return nil
	}

	if tm, ok := textMarshalerFor(v); ok {
		b, err := tm.MarshalText()
		if err != nil {
			return err
		}
		s := string(b)
		*out = append(*out, Field{Key: key, Value: &s})
		return nil
	}

	switch v.Kind() {
	case reflect.Struct:
		if v.Type() == nodeType {
			n := v.Interface().(Node)
			nCopy := n
			nCopy.Name = key
			*out = append(*out, Field{Key: key, Child: &nCopy})
			return nil
		}
		fields, err := encodeStructFields(v)
		if err != nil {
			return err
		}
		*out = append(*out, Field{Key: key, Child: &Node{Name: key, Fields: fields}})
		return nil

	case reflect.Map:
		fields, err := encodeMapFields(v)
		if err != nil {
			return err
		}
		*out = append(*out, Field{Key: key, Child: &Node{Name: key, Fields: fields}})
		return nil

	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
			s := string(v.Bytes())
			*out = append(*out, Field{Key: key, Value: &s})
			return nil
		}
		if tuple {
			parts := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				s, err := scalarString(v.Index(i))
				if err != nil {
					return err
				}
				parts[i] = s
			}
			s := formatTuple(parts)
			*out = append(*out, Field{Key: key, Value: &s})
			return nil
		}
		for i := 0; i < v.Len(); i++ {
			if err := appendEncodedField(out, key, v.Index(i), false); err != nil {
				return err
			}
		}
		return nil

	default:
		s, err := scalarString(v)
		if err != nil {
			return err
		}
		*out = append(*out, Field{Key: key, Value: &s})
		return nil
	}
}

// scalarString renders a basic-kinded Value as scalar text.
func scalarString(v reflect.Value) (string, error) {
	if tm, ok := textMarshalerFor(v); ok {
		b, err := tm.MarshalText()
		return string(b), err
	}
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return formatBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'g', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), nil
	default:
		return "", &UnsupportedTypeError{Type: v.Type()}
	}
}

func encodeStructFields(v reflect.Value) ([]Field, error) {
	infos := structFieldInfos(v.Type())
	var out []Field
	for _, fi := range infos {
		fv, ok := fieldByIndexSafe(v, fi.index)
		if !ok {
			continue // nil embedded pointer chain: nothing to encode
		}
		if fi.tag.omitempty && isEmptyValue(fv) {
			continue
		}
		if err := appendEncodedField(&out, fi.tag.name, fv, fi.tag.tuple); err != nil {
			return nil, fmt.Errorf("field %q: %w", fi.tag.name, err)
		}
	}
	return out, nil
}

func fieldByIndexSafe(v reflect.Value, index []int) (reflect.Value, bool) {
	for i, x := range index {
		if i > 0 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					return reflect.Value{}, false
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v, true
}

func encodeMapFields(v reflect.Value) ([]Field, error) {
	keys := v.MapKeys()
	strKeys := make([]string, len(keys))
	for i, k := range keys {
		strKeys[i] = fmt.Sprint(k.Interface())
	}
	sort.Strings(strKeys)

	var out []Field
	for _, sk := range strKeys {
		mv := v.MapIndex(reflect.ValueOf(sk).Convert(v.Type().Key()))
		if mv.Kind() == reflect.Slice && mv.Type().Elem().Kind() != reflect.Uint8 {
			for i := 0; i < mv.Len(); i++ {
				if err := appendEncodedField(&out, sk, mv.Index(i), false); err != nil {
					return nil, err
				}
			}
			continue
		}
		if err := appendEncodedField(&out, sk, mv, false); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// isEmptyValue mirrors encoding/json's isEmptyValue.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// writeFields serializes a Field slice to buf using the given prefix and
// per-level indent.
func writeFields(buf *bytes.Buffer, fields []Field, prefix, indent string) {
	for _, f := range fields {
		if f.Child != nil {
			buf.WriteString(prefix)
			buf.WriteString(f.Key)
			buf.WriteByte('\n')
			buf.WriteString(prefix)
			buf.WriteString("{\n")
			writeFields(buf, f.Child.Fields, prefix+indent, indent)
			buf.WriteString(prefix)
			buf.WriteString("}\n")
			continue
		}
		buf.WriteString(prefix)
		buf.WriteString(f.Key)
		buf.WriteString(" = ")
		if f.Value != nil {
			buf.WriteString(*f.Value)
		}
		buf.WriteByte('\n')
	}
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
