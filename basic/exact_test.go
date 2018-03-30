package basic

import (
	"testing"
)

func TestBool_Exact(t *testing.T) {
	TestExactBool(t, ExactEqualer{})
}

func TestInt_Exact(t *testing.T) {
	TestExactInt64(t, ExactEqualer{})
}

func TestFloat64_Exact(t *testing.T) {
	TestExactFloat64(t, ExactEqualer{})
}

func TestString_Exact(t *testing.T) {
	TestExactString(t, ExactEqualer{})
}
