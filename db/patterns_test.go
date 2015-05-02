package db

import "testing"

func TestPatterns(t *testing.T) {
	type pattern interface {
		Matches(string) bool
	}
	var (
		p   *PathPattern
		d   *DomainPattern
		err error
	)
	parsePath := func(s string) {
		p, err = ParsePath(s)
		if err != nil {
			t.Fatalf("%#v failed to parse - %v\n", s, err)
		}
		t.Logf("parsed path pattern %#v\n", s)
	}
	parseDomain := func(s string) {
		d, err = ParseDomain(s)
		if err != nil {
			t.Fatalf("%#v failed to parse - %v\n", s, err)
		}
		t.Logf("parsed domain pattern %#v\n", s)
	}
	expectMatch := func(p pattern, s string) {
		if !p.Matches(s) {
			t.Errorf("pattern match failed - %#v", s)
		}
	}
	expectFail := func(p pattern, s string) {
		if p.Matches(s) {
			t.Errorf("invalid pattern match - %#v", s)
		}
	}

	// paths
	parsePath("")
	expectMatch(p, "")
	expectMatch(p, "/")
	expectMatch(p, "/a")

	parsePath("$")
	expectMatch(p, "")
	expectFail(p, "/")
	expectFail(p, "/a")

	parsePath("/$")
	expectMatch(p, "/")
	expectFail(p, "/a")

	parsePath("/test/")
	expectMatch(p, "/test/1")
	expectFail(p, "/test")

	parsePath("/test/$")
	expectMatch(p, "/test/")
	expectFail(p, "/test")
	expectFail(p, "/test/a")

	parsePath("/test/*/view")
	expectMatch(p, "/test/a/view")
	expectFail(p, "/test/a/1/view")

	parsePath("/te*")
	expectMatch(p, "/test")
	expectFail(p, "/best")

	parsePath("/*st")
	expectMatch(p, "/test")
	expectFail(p, "/teak")

	parsePath("/*a*")
	expectMatch(p, "/fail")
	expectFail(p, "/fill")

	// domain
	parseDomain("example.com")
	expectMatch(d, "example.com")
	expectFail(d, "x.example.com")
	expectFail(d, "something.com")
	expectFail(d, "example.org")

	parseDomain("*.example.com")
	expectMatch(d, "abc.example.com")
	expectFail(d, "example.com")

	parseDomain("example.*")
	expectMatch(d, "example.com")
	expectMatch(d, "example.org")
}
