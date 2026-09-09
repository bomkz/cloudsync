package vtscfg

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal parses config data and stores the result in the value
// pointed to by v, the same way json.Unmarshal does for JSON. v must be a
// non-nil pointer.
//
// The synthetic top-level Node's Fields (i.e. the file's top-level
// blocks) are decoded into v exactly like the Fields of any other Node,
// so a top-level struct field tagged `vts:"PILOTS"` receives the file's
// "PILOTS { ... }" block.
func Unmarshal(data []byte, v interface{}) error {
	root, err := Parse(data)
	if err != nil {
		return err
	}
	return decodeInto(root, v)
}

func decodeInto(root *Node, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}
	return decodeNodeInto(root, rv)
}

// Decoder reads and decodes a single config document from an input
// stream, mirroring json.Decoder's constructor/method shape. Unlike JSON,
// this format has no notion of multiple concatenated top-level documents
// in one stream, so a second call to Decode returns io.EOF.
type Decoder struct {
	r    io.Reader
	done bool
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the entire remaining input from the Decoder's reader and
// stores the result in the value pointed to by v.
func (d *Decoder) Decode(v interface{}) error {
	if d.done {
		return io.EOF
	}
	data, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}
	d.done = true
	if len(bytes.TrimSpace(data)) == 0 {
		return io.EOF
	}
	return Unmarshal(data, v)
}

var nodeType = reflect.TypeOf(Node{})

// indirect walks a chain of pointers, allocating as needed, and returns
// the first non-pointer, settable Value. It leaves interface values
// alone (those are handled by their callers).
func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

// fieldByIndexAlloc walks a VisibleFields-style index path from v,
// allocating any nil embedded pointer structs it passes through.
func fieldByIndexAlloc(v reflect.Value, index []int) reflect.Value {
	for i, x := range index {
		if i > 0 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}

// groupFields groups a Node's Fields by key, preserving first-seen order.
func groupFields(n *Node) (order []string, grouped map[string][]Field) {
	grouped = make(map[string][]Field, len(n.Fields))
	for _, f := range n.Fields {
		if _, ok := grouped[f.Key]; !ok {
			order = append(order, f.Key)
		}
		grouped[f.Key] = append(grouped[f.Key], f)
	}
	return order, grouped
}

// decodeField decodes a single Field (scalar or nested node) into target.
func decodeField(f Field, target reflect.Value) error {
	target = indirect(target)

	if f.Child != nil {
		if target.CanAddr() {
			if u, ok := target.Addr().Interface().(Unmarshaler); ok {
				return u.UnmarshalVTS(f.Child)
			}
		}
		return decodeNodeInto(f.Child, target)
	}

	raw := ""
	if f.Value != nil {
		raw = *f.Value
	}
	return decodeScalarInto(raw, target)
}

// decodeNodeInto decodes a whole Node (struct/map/interface{}/Node
// target) into target, which must already be indirected past pointers by
// the caller (decodeInto passes the pointer itself, which indirect()
// inside handles).
func decodeNodeInto(node *Node, target reflect.Value) error {
	target = indirect(target)

	if target.CanAddr() {
		if u, ok := target.Addr().Interface().(Unmarshaler); ok {
			return u.UnmarshalVTS(node)
		}
	}

	if target.Type() == nodeType {
		target.Set(reflect.ValueOf(*node))
		return nil
	}

	switch target.Kind() {
	case reflect.Interface:
		if target.NumMethod() != 0 {
			return &UnmarshalTypeError{Value: "block", Type: target.Type()}
		}
		m, err := nodeToInterface(node)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(m))
		return nil

	case reflect.Struct:
		return decodeStruct(node, target)

	case reflect.Map:
		return decodeMap(node, target)

	default:
		return &UnmarshalTypeError{Value: "block", Type: target.Type()}
	}
}

type fieldInfo struct {
	index []int
	tag   tagOptions
}

func structFieldInfos(t reflect.Type) []fieldInfo {
	var infos []fieldInfo
	for _, sf := range reflect.VisibleFields(t) {
		if sf.PkgPath != "" {
			continue // unexported
		}
		if sf.Anonymous {
			if _, hasTag := sf.Tag.Lookup("vts"); !hasTag {
				// Its promoted fields already appear as separate entries
				// from reflect.VisibleFields; don't also treat the
				// embedded struct itself as a field to match.
				continue
			}
		}
		raw, hasTag := sf.Tag.Lookup("vts")
		opt := parseTag(raw, hasTag, sf.Name)
		if opt.skip {
			continue
		}
		infos = append(infos, fieldInfo{index: append([]int(nil), sf.Index...), tag: opt})
	}
	return infos
}

func decodeStruct(node *Node, target reflect.Value) error {
	infos := structFieldInfos(target.Type())
	_, grouped := groupFields(node)

	for _, fi := range infos {
		matches, ok := grouped[fi.tag.name]
		if !ok {
			// Fall back to a case-insensitive match, same as encoding/json
			// does for untagged fields.
			for k, v := range grouped {
				if strings.EqualFold(k, fi.tag.name) {
					matches, ok = v, true
					break
				}
			}
		}
		if !ok {
			continue
		}
		fv := fieldByIndexAlloc(target, fi.index)
		if err := assignMatches(matches, fv); err != nil {
			return fmt.Errorf("field %q: %w", fi.tag.name, err)
		}
	}
	return nil
}

func decodeMap(node *Node, target reflect.Value) error {
	t := target.Type()
	if t.Key().Kind() != reflect.String {
		return &UnmarshalTypeError{Value: "block", Type: t}
	}
	if target.IsNil() {
		target.Set(reflect.MakeMap(t))
	}
	order, grouped := groupFields(node)
	elemType := t.Elem()
	for _, key := range order {
		matches := grouped[key]
		elem := reflect.New(elemType).Elem()
		if err := assignMatches(matches, elem); err != nil {
			return fmt.Errorf("key %q: %w", key, err)
		}
		target.SetMapIndex(reflect.ValueOf(key).Convert(t.Key()), elem)
	}
	return nil
}

// assignMatches assigns every Field sharing a key onto target, which may
// be a slice/array (one element per match, or a single tuple-shaped
// scalar match unpacked element-wise) or a scalar/struct/map (in which
// case the *last* match wins, mirroring encoding/json's duplicate-key
// behavior).
func assignMatches(matches []Field, target reflect.Value) error {
	target = indirect(target)

	switch target.Kind() {
	case reflect.Slice:
		if target.Type().Elem().Kind() == reflect.Uint8 {
			if len(matches) > 0 && matches[len(matches)-1].Value != nil {
				target.SetBytes([]byte(*matches[len(matches)-1].Value))
			}
			return nil
		}
		if len(matches) == 1 && matches[0].Value != nil && isTuple(*matches[0].Value) {
			return decodeScalarInto(*matches[0].Value, target)
		}
		sl := reflect.MakeSlice(target.Type(), len(matches), len(matches))
		for i, m := range matches {
			if err := decodeField(m, sl.Index(i)); err != nil {
				return err
			}
		}
		target.Set(sl)
		return nil

	case reflect.Array:
		if len(matches) == 1 && matches[0].Value != nil && isTuple(*matches[0].Value) {
			return decodeScalarInto(*matches[0].Value, target)
		}
		for i := 0; i < target.Len() && i < len(matches); i++ {
			if err := decodeField(matches[i], target.Index(i)); err != nil {
				return err
			}
		}
		return nil

	default:
		return decodeField(matches[len(matches)-1], target)
	}
}

// decodeScalarInto decodes raw scalar text into target (already
// indirected).
func decodeScalarInto(raw string, target reflect.Value) error {
	target = indirect(target)

	if target.CanAddr() {
		if u, ok := target.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText([]byte(raw))
		}
	}

	switch target.Kind() {
	case reflect.Interface:
		if target.NumMethod() != 0 {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		target.Set(reflect.ValueOf(scalarToInterface(raw)))
		return nil

	case reflect.String:
		target.SetString(raw)
		return nil

	case reflect.Bool:
		if raw == "" {
			return nil
		}
		b, ok := parseBool(raw)
		if !ok {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		target.SetBool(b)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if raw == "" {
			return nil
		}
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		if target.OverflowInt(n) {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		target.SetInt(n)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if raw == "" {
			return nil
		}
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		if target.OverflowUint(n) {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		target.SetUint(n)
		return nil

	case reflect.Float32, reflect.Float64:
		if raw == "" {
			return nil
		}
		f, err := strconv.ParseFloat(raw, target.Type().Bits())
		if err != nil {
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		target.SetFloat(f)
		return nil

	case reflect.Slice:
		if !isTuple(raw) {
			if raw == "" {
				return nil
			}
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		parts := splitTuple(raw)
		sl := reflect.MakeSlice(target.Type(), len(parts), len(parts))
		for i, p := range parts {
			if err := decodeScalarInto(p, sl.Index(i)); err != nil {
				return err
			}
		}
		target.Set(sl)
		return nil

	case reflect.Array:
		if !isTuple(raw) {
			if raw == "" {
				return nil
			}
			return &UnmarshalTypeError{Value: raw, Type: target.Type()}
		}
		parts := splitTuple(raw)
		for i := 0; i < target.Len() && i < len(parts); i++ {
			if err := decodeScalarInto(parts[i], target.Index(i)); err != nil {
				return err
			}
		}
		return nil

	default:
		return &UnmarshalTypeError{Value: raw, Type: target.Type()}
	}
}

// scalarToInterface converts raw scalar text into a generic Go value the
// way encoding/json converts a JSON scalar for an interface{} target:
// bool, float64, or string (plus []interface{} for our tuple syntax,
// which has no JSON equivalent).
func scalarToInterface(raw string) interface{} {
	if raw == "" {
		return ""
	}
	if b, ok := parseBool(raw); ok {
		return b
	}
	if isTuple(raw) {
		parts := splitTuple(raw)
		arr := make([]interface{}, len(parts))
		for i, p := range parts {
			arr[i] = scalarToInterface(p)
		}
		return arr
	}
	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		return f
	}
	return raw
}

// nodeToInterface converts a whole Node into a generic
// map[string]interface{} tree, promoting repeated keys to []interface{}.
func nodeToInterface(node *Node) (interface{}, error) {
	order, grouped := groupFields(node)
	out := make(map[string]interface{}, len(order))
	for _, key := range order {
		matches := grouped[key]
		if len(matches) == 1 {
			v, err := fieldToInterface(matches[0])
			if err != nil {
				return nil, err
			}
			out[key] = v
			continue
		}
		arr := make([]interface{}, len(matches))
		for i, m := range matches {
			v, err := fieldToInterface(m)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		out[key] = arr
	}
	return out, nil
}

func fieldToInterface(f Field) (interface{}, error) {
	if f.Child != nil {
		return nodeToInterface(f.Child)
	}
	raw := ""
	if f.Value != nil {
		raw = *f.Value
	}
	return scalarToInterface(raw), nil
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
