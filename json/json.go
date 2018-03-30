// Package json provides functionality for comparing JSON strings.
package json

import (
	"encoding/json"
	"reflect"

	"github.com/haimberger/compare/basic"
)

// Equaler implements functions for determining if JSON strings are equal.
type Equaler struct {
	Basic basic.Equaler
}

// Equal determines if two JSON strings represent the same value.
func (e Equaler) Equal(s1, s2 []byte) (bool, error) {
	var v1, v2 interface{}

	if err := json.Unmarshal(s1, &v1); err != nil {
		return false, err
	}
	if err := json.Unmarshal(s2, &v2); err != nil {
		return false, err
	}

	return e.equal(v1, v2), nil
}

func (e Equaler) equal(v1, v2 interface{}) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}

	if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
		return false
	}

	switch a := v1.(type) {
	case []interface{}:
		return e.equalArrays(a, v2.([]interface{}))
	case map[string]interface{}:
		return e.equalObjects(a, v2.(map[string]interface{}))
	case bool:
		return e.Basic.Bool(a, v2.(bool))
	case float64:
		return e.Basic.Float64(a, v2.(float64))
	case string:
		return e.Basic.String(a, v2.(string))
	default:
		// should never happen (https://golang.org/pkg/encoding/json/#Unmarshal)
		return reflect.DeepEqual(v1, v2)
	}
}

func (e Equaler) equalArrays(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if !e.equal(v, b[i]) {
			return false
		}
	}

	return true
}

func (e Equaler) equalObjects(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if !e.equal(v, b[k]) {
			return false
		}
	}

	return true
}
