package models

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

const testDbPath = "../out/test.db"

func failIf(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	db, err := bolt.Open(testDbPath, 0660, nil)
	panicIf(err)
	err = InitDb(db)
	panicIf(err)
	r := m.Run()
	db.Close()
	os.Exit(r)
}

func TestGroups(t *testing.T) {
	g := &Group{
		Name:         "testGroup",
		Proto:        ProtoBoth,
		System:       "HC128",
		SkipFragment: true,
	}
	err := AddGroup(g)
	failIf(t, err)
	ret, err := GetGroup(g.Key)
	failIf(t, err)
	if *ret != *g {
		t.Error("retrieved value not same as inserted")
		t.Logf("%+v != %+v\n", ret, g)
	}
}
