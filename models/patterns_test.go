package models

import "testing"

func TestPathPattern(t *testing.T) {
	var (
		p   *PathPattern
		err error
	)
	parse := func(s string) {
		p, err = ParsePath(s)
		if err != nil {
			t.Fatalf("%#v failed to parse - %v\n", s, err)
		}
		t.Logf("parsed pattern %#v\n", s)
	}
	expectMatch := func(s string) {
		if !p.Matches(s) {
			t.Errorf("pattern match failed - %#v", s)
		}
	}
	expectFail := func(s string) {
		if p.Matches(s) {
			t.Errorf("invalid pattern match - %#v", s)
		}
	}
	parse("")
	expectMatch("")
	expectMatch("/")
	expectMatch("/a")
	parse("$")
	expectMatch("")
	expectFail("/")
	expectFail("/a")
	parse("/$")
	expectMatch("/")
	expectFail("/a")
	parse("/test/")
	expectMatch("/test/1")
	expectFail("/test")
	parse("/test/$")
	expectMatch("/test/")
	expectFail("/test")
	expectFail("/test/a")
	parse("/test/*/view")
	expectMatch("/test/a/view")
	expectFail("/test/a/1/view")
	parse("/te*")
	expectMatch("/test")
	expectFail("/best")
	parse("/*st")
	expectMatch("/test")
	expectFail("/teak")
	parse("/*a*")
	expectMatch("/fail")
	expectFail("/fill")
}
