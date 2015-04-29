package models

import (
	"errors"
	"math/rand"
	"net/url"

	"github.com/boltdb/bolt"
)

var (
	ErrNoSuchGroupKey = errors.New("non-existent group key")
	ErrInvURL         = errors.New("invalid URL")
	ErrInvProto       = errors.New("wrong protocol")
	ErrNoDomainMatch  = errors.New("URL doesn't match any domain pattern")
	ErrNoPathMatch    = errors.New("URL doesn't match any path pattern")
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type protoSpec int

const (
	ProtoPlain protoSpec = iota + 1
	ProtoSecure
	ProtoBoth
)

type Group struct {
	Key          string
	Name         string
	Proto        protoSpec
	System       string
	SkipFragment bool
}

func (g *Group) IsValid(loc string) (string, error) {
	u, err := url.Parse(loc)
	if err != nil {
		return "", ErrInvURL
	}
	protoValid := false
	if (g.Proto&ProtoPlain) == 1 && u.Scheme == "http" {
		protoValid = true
	}
	if (g.Proto&ProtoSecure) == 1 && u.Scheme == "https" {
		protoValid = true
	}
	if !protoValid {
		return "", ErrInvProto
	}
	if g.SkipFragment {
		u.Fragment = ""
	}
	//TODO: verify path and domain
	return u.String(), nil
}

const (
	chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	groupKeyLen = 6
)

func genGroupKey() string {
	id := make([]byte, groupKeyLen)
	for i := range id {
		r := rand.Intn(len(chars))
		id[i] = chars[r]
	}
	return string(id)
}

func keyExists(key string) bool {
	//TODO
	var exists bool
	db.View(func(tx *bolt.Tx) error {
		groupBucket := tx.Bucket(groupBucKey)
		v := groupBucket.Get([]byte(key))
		if v == nil {
			exists = false
		} else {
			exists = true
		}
		return nil
	})
	return exists
}

func UpdateGroup(group *Group) error {
	err := db.Update(func(tx *bolt.Tx) error {
		groupBucket := tx.Bucket(groupBucKey)
		v, err := encGob(group)
		if err != nil {
			return err
		}
		return groupBucket.Put([]byte(group.Key), v)
	})
	return err
}

func GetGroup(key string) (*Group, error) {
	g := new(Group)
	err := db.View(func(tx *bolt.Tx) error {
		groupBucket := tx.Bucket(groupBucKey)
		v := groupBucket.Get([]byte(key))
		if v == nil {
			return ErrNoSuchGroupKey
		}
		return decGob(v, g)
	})
	if err != nil {
		return nil, err
	}
	return g, nil
}

func AddGroup(group *Group) error {
	key := genGroupKey()
	for keyExists(key) {
		key = genGroupKey()
	}
	group.Key = key
	return UpdateGroup(group)
}
