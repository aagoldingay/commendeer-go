package utilities

import "testing"

func Test_generateCode_CorrectLength(t *testing.T) {
	l := 10
	c := generateCode(l)
	if len(c) != l {
		t.Errorf("generatecode failed: expected %v, actual %v", l, len(c))
	}
}
