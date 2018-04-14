package compare

import (
	"encoding/json"
	"reflect"
	"regexp"
	"sort"

	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

// JSONDiff represents the differences between two JSON values.
type JSONDiff struct {
	left map[string]interface{}
	ds   []gojsondiff.Delta
}

// Deltas returns Deltas that describe individual differences between two JSON values.
func (d *JSONDiff) Deltas() []gojsondiff.Delta {
	return d.ds
}

// Modified returns true if JSONDiff has at least one Delta.
func (d *JSONDiff) Modified() bool {
	return len(d.ds) > 0
}

// Format returns a string representation of the differences between two JSON values.
func (d *JSONDiff) Format(coloring bool) (string, error) {
	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       coloring,
	}
	af := formatter.NewAsciiFormatter(d.left, config)
	diff, err := af.Format(d)
	if err != nil {
		return "", err
	}

	// remove wrapping object for values of basic types
	basicRE, err := regexp.Compile(`^\s*{\n-\s*"\$":\s*([\s\S]*)\+\s*"\$":\s*([\s\S]*)\n\s*}`)
	if err != nil {
		return "", err
	}
	unwrapped := basicRE.ReplaceAllString(diff, "- $1+ $2")

	// remove wrapping object for arrays and objects
	nonBasicRE, err := regexp.Compile(`^\s*{\n\s*"\$":([\s\S]*)\n\s}`)
	if err != nil {
		return "", err
	}
	unwrapped = nonBasicRE.ReplaceAllString(unwrapped, "$1")

	// decrease indentation for remaining lines
	indentRE, err := regexp.Compile(`(\n.)\s{2}`)
	if err != nil {
		return "", err
	}
	return indentRE.ReplaceAllString(unwrapped, "$1"), nil
}

// JSONDiffer compares JSON strings.
type JSONDiffer struct {
	// Basic specifies how values of basic types should be compared.
	Basic BasicEqualer
}

// Equal determines if two JSON strings represent the same value.
// Returns an error iff the strings don't adhere to the JSON syntax.
func (jd JSONDiffer) Equal(left, right []byte) (bool, error) {
	d, err := jd.Compare(left, right)
	if err != nil {
		return false, err
	}
	return !d.Modified(), nil
}

// Compare returns the differences between two JSON strings.
// Returns an error iff the strings don't adhere to the JSON syntax.
func (jd JSONDiffer) Compare(left, right []byte) (*JSONDiff, error) {
	var l, r interface{}

	if err := json.Unmarshal(left, &l); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(right, &r); err != nil {
		return nil, err
	}

	// add explicit root in case the values are arrays or plain values (not objects)
	leftMap := map[string]interface{}{"$": l}
	rightMap := map[string]interface{}{"$": r}
	d := jd.compareMaps(leftMap, rightMap)
	d.left = leftMap
	return d, nil
}

// cf. https://github.com/yudai/gojsondiff/blob/master/gojsondiff.go#L66-L74
func (jd JSONDiffer) compareMaps(left, right map[string]interface{}) *JSONDiff {
	ds := jd.mapDeltas(left, right)
	return &JSONDiff{ds: ds}
}

// cf. https://github.com/yudai/gojsondiff/blob/master/gojsondiff.go#L235-L279
func (jd JSONDiffer) compare(pos gojsondiff.Position, left, right interface{}) (bool, gojsondiff.Delta) {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false, gojsondiff.NewModified(pos, left, right)
	}

	switch l := left.(type) {
	case []interface{}:
		if ds := jd.sliceDeltas(l, right.([]interface{})); len(ds) > 0 {
			return false, gojsondiff.NewArray(pos, ds)
		}
	case map[string]interface{}:
		if ds := jd.mapDeltas(l, right.(map[string]interface{})); len(ds) > 0 {
			return false, gojsondiff.NewObject(pos, ds)
		}
	default:
		return jd.valueDelta(pos, left, right)
	}

	return true, nil
}

// cf. https://github.com/yudai/gojsondiff/blob/master/gojsondiff.go#L125-L233
// Note that this implementation is much more primitive. There's no attempt to
// find a longest common sequence and base differences on that. We just compare
// values index by index.
func (jd JSONDiffer) sliceDeltas(left, right []interface{}) []gojsondiff.Delta {
	var ds []gojsondiff.Delta

	for i, leftVal := range left {
		if i < len(right) {
			if same, d := jd.compare(gojsondiff.Index(i), leftVal, right[i]); !same {
				ds = append(ds, d)
			}
		} else {
			ds = append(ds, gojsondiff.NewDeleted(gojsondiff.Index(i), leftVal))
		}
	}

	for i := len(left); i < len(right); i++ {
		ds = append(ds, gojsondiff.NewAdded(gojsondiff.Index(i), right[i]))
	}

	return ds
}

// cf. https://github.com/yudai/gojsondiff/blob/master/gojsondiff.go#L86-L112
func (jd JSONDiffer) mapDeltas(left, right map[string]interface{}) []gojsondiff.Delta {
	var ds []gojsondiff.Delta

	keys := sortedKeys(left) // stabilize delta order
	for _, key := range keys {
		if rightVal, ok := right[key]; ok {
			if same, d := jd.compare(gojsondiff.Name(key), left[key], rightVal); !same {
				ds = append(ds, d)
			}
		} else {
			ds = append(ds, gojsondiff.NewDeleted(gojsondiff.Name(key), left[key]))
		}
	}

	keys = sortedKeys(right) // stabilize delta order
	for _, key := range keys {
		if _, ok := left[key]; !ok {
			ds = append(ds, gojsondiff.NewAdded(gojsondiff.Name(key), right[key]))
		}
	}

	return ds
}

// valueDelta returns the Delta (if any) for two basic values (null, boolean, number, string).
// Rather than just using reflect.DeepEqual(), as gojsondiff does, we use a custom BasicEqualer.
func (jd *JSONDiffer) valueDelta(pos gojsondiff.Position, left, right interface{}) (bool, gojsondiff.Delta) {
	var same bool

	switch l := left.(type) {
	case nil:
		same = left == right
	case bool:
		same = jd.Basic.Bool(l, right.(bool))
	case float64:
		same = jd.Basic.Float64(l, right.(float64))
	case string:
		same = jd.Basic.String(l, right.(string))
	default:
		// should never happen (https://golang.org/pkg/encoding/json/#Unmarshal)
		same = reflect.DeepEqual(left, right)
	}

	if !same {
		return false, gojsondiff.NewModified(pos, left, right)
	}

	return true, nil
}

// cf. https://github.com/yudai/gojsondiff/blob/master/gojsondiff.go#L409-L416
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
