// Package json provides functionality for comparing JSON strings.
package json

import (
	"encoding/json"
	"reflect"

	"github.com/haimberger/compare/basic"
)

// Equal determines if two JSON strings represent the same value.
func Equal(s1, s2 []byte, be basic.Equaler) (bool, error) {
	var v1, v2 interface{}

	if err := json.Unmarshal(s1, &v1); err != nil {
		return false, err
	}
	if err := json.Unmarshal(s2, &v2); err != nil {
		return false, err
	}

	return equal(v1, v2, be), nil
}

func equal(v1, v2 interface{}, be basic.Equaler) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}

	if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
		return false
	}

	switch a := v1.(type) {
	case []interface{}:
		return equalArrays(a, v2.([]interface{}), be)
	case map[string]interface{}:
		return equalObjects(a, v2.(map[string]interface{}), be)
	case bool:
		return be.Bool(a, v2.(bool))
	case float64:
		return be.Float64(a, v2.(float64))
	case string:
		return be.String(a, v2.(string))
	default:
		// should never happen (https://golang.org/pkg/encoding/json/#Unmarshal)
		return reflect.DeepEqual(v1, v2)
	}
}

func equalArrays(a, b []interface{}, be basic.Equaler) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if !equal(v, b[i], be) {
			return false
		}
	}

	return true
}

func equalObjects(a, b map[string]interface{}, be basic.Equaler) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if !equal(v, b[k], be) {
			return false
		}
	}

	return true
}
