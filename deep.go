package compare

import (
	"fmt"
	"reflect"
)

// DeepEqualer provides functions for deeep comparison of values.
type DeepEqualer struct {
	// Basic specifies how values of basic types should be compared.
	Basic BasicEqualer
}

// Equal determines if two values contain the same information.
//
// Shortcomings:
//
// 	- we assume that there are no loops
// 	- there are no plans to support the comparison of functions
// 	- some basic types, like uint and complex128, aren't yet supported
// 	- behaviour when comparing values with an empty struct (e.g. `type self struct{}`)
// 	  is different from the behaviour of reflect.DeepEqual(); not quite sure why...
//
// The implementation closely follows the implementation of reflect.DeepEqual().
func (e DeepEqualer) Equal(a, b interface{}) (bool, error) {
	return e.equal(reflect.ValueOf(a), reflect.ValueOf(b))
}

// nolint: gocyclo
// The complexity is currently 15 (greater than the desired maximum of 10).
// I disabled the gocyclo check, because I can't think of a way to reduce the
// cyclomatic complexity in a way that really feels like an improvement.
// In any case, I think that the code is easy to follow as it is, and the test
// coverage for this function is 100% despite the high number of execution paths.
func (e DeepEqualer) equal(v1, v2 reflect.Value) (bool, error) {
	if !v1.IsValid() || !v2.IsValid() { // at least one underlying value was nil
		return v1.IsValid() == v2.IsValid(), nil
	}

	if v1.Type() != v2.Type() {
		return false, nil
	}

	switch v1.Kind() {
	case reflect.Array:
		return e.equalArrays(v1, v2)
	case reflect.Bool:
		return e.Basic.Bool(v1.Bool(), v2.Bool()), nil
	case reflect.Float32, reflect.Float64:
		return e.Basic.Float64(v1.Float(), v2.Float()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.Basic.Int64(v1.Int(), v2.Int()), nil
	case reflect.Interface:
		return e.equalInterfaces(v1, v2)
	case reflect.Map:
		return e.equalMaps(v1, v2)
	case reflect.Ptr:
		return e.equalPointers(v1, v2)
	case reflect.Slice:
		return e.equalSlices(v1, v2)
	case reflect.String:
		return e.Basic.String(v1.String(), v2.String()), nil
	case reflect.Struct:
		return e.equalStructs(v1, v2)
	default: // unsupported or incomparable type (e.g. Complex128, Func)
		return false, fmt.Errorf("type %v not supported", v1.Type())
	}
}

func (e DeepEqualer) equalArrays(v1, v2 reflect.Value) (bool, error) {
	for i := 0; i < v1.Len(); i++ {
		if eq, err := e.equal(v1.Index(i), v2.Index(i)); err != nil || !eq {
			return false, err
		}
	}
	return true, nil
}

func (e DeepEqualer) equalInterfaces(v1, v2 reflect.Value) (bool, error) {
	if v1.IsNil() || v2.IsNil() {
		return v1.IsNil() == v2.IsNil(), nil
	}
	return e.equal(v1.Elem(), v2.Elem())
}

func (e DeepEqualer) equalMaps(v1, v2 reflect.Value) (bool, error) {
	if v1.IsNil() != v2.IsNil() {
		return false, nil
	}
	if v1.Len() != v2.Len() {
		return false, nil
	}
	if v1.Pointer() == v2.Pointer() {
		return true, nil
	}
	for _, k := range v1.MapKeys() {
		val1, val2 := v1.MapIndex(k), v2.MapIndex(k)
		if !val1.IsValid() || !val2.IsValid() {
			return false, nil
		}
		if eq, err := e.equal(val1, val2); err != nil || !eq {
			return false, err
		}
	}
	return true, nil
}

func (e DeepEqualer) equalPointers(v1, v2 reflect.Value) (bool, error) {
	if v1.Pointer() == v2.Pointer() {
		return true, nil
	}
	return e.equal(v1.Elem(), v2.Elem())
}

func (e DeepEqualer) equalSlices(v1, v2 reflect.Value) (bool, error) {
	if v1.IsNil() != v2.IsNil() {
		return false, nil
	}
	if v1.Len() != v2.Len() {
		return false, nil
	}
	if v1.Pointer() == v2.Pointer() {
		return true, nil
	}
	for i := 0; i < v1.Len(); i++ {
		if eq, err := e.equal(v1.Index(i), v2.Index(i)); err != nil || !eq {
			return false, err
		}
	}
	return true, nil
}

func (e DeepEqualer) equalStructs(v1, v2 reflect.Value) (bool, error) {
	for i := 0; i < v1.NumField(); i++ {
		if eq, err := e.equal(v1.Field(i), v2.Field(i)); err != nil || !eq {
			return false, err
		}
	}
	return true, nil
}