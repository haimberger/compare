package compare

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func ExampleJSONDiffer_Compare() {
	jd := JSONDiffer{TolerantBasicEqualer{Float64Tolerance: 0.1}}
	d, err := jd.Compare(
		[]byte(`{"x": 1.6, "y": [3.8, "hello"]}`),
		[]byte(`{"x": 1.57, "y": [3.6, "hello"], "z": 0}`))
	if err != nil {
		fmt.Println(err)
		return
	}
	diff, err := d.Format(false) // no colouring
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(diff)
	// Output:
	//  {
	//    "x": 1.6,
	//    "y": [
	// -    0: 3.8,
	// +    0: 3.6,
	//      1: "hello"
	//    ]
	// +  "z": 0
	//  }
}

func TestJSONDiffer_Equal_exact(t *testing.T) {
	type testCase struct {
		a        string
		b        string
		expected bool
		err      string
	}
	tcs := []testCase{
		{``, ``, false, "unexpected end of JSON input"},
		{`""`, ``, false, "unexpected end of JSON input"},
		{`undefined`, `undefined`, false, "invalid character 'u' looking for beginning of value"},
		{`null`, `null`, true, ""},
		{`null`, `false`, false, ""},
		{`null`, `""`, false, ""},
		{`false`, `false`, true, ""},
		{`false`, `true`, false, ""},
		{`1`, `1`, true, ""},
		{`1`, `2`, false, ""},
		{`0.1`, `0.1`, true, ""},
		{`0.1`, `0.2`, false, ""},
		{`"foo"`, `"foo"`, true, ""},
		{`"foo"`, `"bar"`, false, ""},
		{`[false, 1, 0.1, "foo"]`, `[false,1,0.100,"foo"]`, true, ""},
		{`[false, 1, 0.1, "foo"]`, `[false, 1, 0.2, "foo"]`, false, ""},
		{`[false, 1, 0.1, "foo"]`, `[false, 0.1, 1, "foo"]`, false, ""},
		{`[false, 1, 0.1]`, `[false, 1, 0.1, "foo"]`, false, ""},
		{`[false, 1, 0.1, "foo"]`, `[false, 1, 0.1]`, false, ""},
		{
			`{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			`{"b": [1, {"d": "foo", "c": 0.1}],"a":false}`,
			true,
			"",
		},
		{
			`{"b": [1, {"c": 0.1, "d": "foo"}]}`,
			`{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			false,
			"",
		},
		{
			`{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			`{"b": [1, {"c": 0.1, "d": "foo"}]}`,
			false,
			"",
		},
		{
			`{"a": false, "b": [1, {"c": 0.1, "d": "foo"}]}`,
			`{"a": false, "b": [1, {"d": 0.1, "c": "foo"}]}`,
			false,
			"",
		},
	}
	// since we specify no tolerances, the equaler will compare values exactly
	e := &JSONDiffer{TolerantBasicEqualer{}}
	for _, tc := range tcs {
		actual, err := e.Equal([]byte(tc.a), []byte(tc.b))
		if err != nil {
			if tc.err == "" {
				t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
			} else if err.Error() != tc.err {
				t.Errorf("[%v == %v] expected error %v; got %v", tc.a, tc.b, tc.err, err)
			}
		} else if tc.err != "" {
			t.Errorf("[%v == %v] expected error %v; got nothing", tc.a, tc.b, tc.err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestJSONDiffer_Equal_tolerant(t *testing.T) {
	type testCase struct {
		a        string
		b        string
		expected bool
	}
	tcs := []testCase{
		{`0.1`, `0.151`, false},
		{`[0.1,0.1,0.1,0.1,0.1]`, `[0.05,0.1,0.10,0.14,0.15]`, true},
		{`"foo_1_1"`, `"foo_2_1"`, false},
		{`["foo_1_1", "foo_1_1"]`, `["foo_1_2", "foo_1_abc"]`, true},
		{`"2018-03-30T16:41:11.509Z"`, `"2018-03-30T16:41:12.510Z"`, false},
		{
			`["2018-03-30T16:41:11.509Z","2018-03-30T16:41:12.5Z"]`,
			`["2018-03-30T16:41:11.509Z","2018-03-30T16:41:12.509Z"]`,
			true,
		},
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	e := &JSONDiffer{TolerantBasicEqualer{
		Float64Tolerance: 0.05,
		// ignore everything after last underscore
		StringTransformer: SubstringDeleter{Regexp: regexp.MustCompile("_[^_]*$")},
		TimeLayout:        time.RFC3339Nano,
		TimeTolerance:     tolerance,
	}}
	for _, tc := range tcs {
		actual, err := e.Equal([]byte(tc.a), []byte(tc.b))
		if err != nil {
			t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestJSONDiffer_Compare(t *testing.T) {
	type testCase struct {
		a        string
		b        string
		expected string
	}
	tcs := []testCase{
		{
			`"hi"`,
			`"hello"`,
			"- \"hi\"\n+ \"hello\"\n",
		},
		{
			`[1, 0.2]`,
			`[1.04, 0.13]`,
			` [
   0: 1,
-  1: 0.2
+  1: 0.13
 ]
`,
		},
		{
			`{"a": false, "b": [1, {"c": 0.2, "d": "foo"}]}`,
			`{"a": false, "b": [1.04, {"c": 0.13, "d": "foo"}]}`,
			` {
   "a": false,
   "b": [
     0: 1,
     1: {
-      "c": 0.2,
+      "c": 0.13,
       "d": "foo"
     }
   ]
 }
`,
		},
	}
	tolerance, err := time.ParseDuration("1s")
	if err != nil {
		t.Fatal(err)
	}
	e := &JSONDiffer{TolerantBasicEqualer{
		Float64Tolerance: 0.05,
		// ignore everything after last underscore
		StringTransformer: SubstringDeleter{Regexp: regexp.MustCompile("_[^_]*$")},
		TimeLayout:        time.RFC3339Nano,
		TimeTolerance:     tolerance,
	}}
	for _, tc := range tcs {
		d, err := e.Compare([]byte(tc.a), []byte(tc.b))
		if err != nil {
			t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
			continue
		}
		actual, err := d.Format(false)
		if err != nil {
			t.Errorf("[%v == %v] %v", tc.a, tc.b, err)
		} else if actual != tc.expected {
			t.Errorf("[%v == %v] expected %v; got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
