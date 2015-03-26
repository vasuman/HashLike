package models

import (
	"bytes"
	"testing"
)

func TestLocationModel(t *testing.T) {
	db, err := InitDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	const loc = "https://www.varav.in/test/post.html"
	l1, err := addOrGetLocation(loc)
	if err != nil {
		t.Fatal(err)
	}
	l2, err := getLocation(loc)
	if err != nil {
		t.Fatal(err)
	}
	if l1.id != l2.id {
		t.Error("ids not equal")
	}
	if !bytes.Equal(l1.hash, l2.hash) {
		t.Error("hashes not equal")
	}
	t.Log("location model working!")
}
