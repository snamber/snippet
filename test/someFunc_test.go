package main

import "testing"

func Test_someFunc(t *testing.T) {

	// START SNIPPET someFunc
	l := someFunc("some input string")
	// l == 17
	// END SNIPPET

	want := 17
	if l != want {
		t.Errorf("somefunc() = %v, want %v", l, want)
	}
}
