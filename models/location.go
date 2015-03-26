package models

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"net/url"
	"regexp"
)

type Location struct {
	id   int64
	URL  string
	Hash []byte
}

type Site struct {
	Domain *regexp.Regexp
	Paths  []*regexp.Regexp
}

var (
	ErrInvalidURL     = errors.New("supplied location is not a valid URL")
	ErrNoSuchLocation = errors.New("no such location in database")
	ErrLocFailMatch   = errors.New("location failed to match allowed sites constraint")
)

var (
	stmtGetLoc *sql.Stmt
	stmtAddLoc *sql.Stmt
)

var allowedSites []*Site

func AddSite(site *Site) {
	allowedSites = append(allowedSites, site)
}

func locHash(loc string) []byte {
	h := sha256.Sum256([]byte(loc))
	return h[:]
}

func GetLocation(loc string) (*Location, error) {
	row := stmtGetLoc.QueryRow(loc)
	l := new(Location)
	err := row.Scan(&l.id, &l.URL, &l.Hash)
	if err == nil {
		return l, nil
	} else if err == sql.ErrNoRows {
		return nil, ErrNoSuchLocation
	}
	return nil, err
}

func checkLocation(loc string) error {
	u, err := url.Parse(loc)
	if err != nil {
		return ErrInvalidURL
	}
	for _, site := range allowedSites {
		if !site.Domain.MatchString(u.Host) {
			continue
		}
		for _, pathExp := range site.Paths {
			if pathExp.MatchString(u.Path) {
				return nil
			}
		}
	}
	return ErrLocFailMatch
}

func AddLocation(loc string) (*Location, error) {
	err := checkLocation(loc)
	if err != nil {
		return nil, err
	}
	// Adding new location
	l := &Location{
		Hash: locHash(loc),
		URL:  loc,
	}
	res, err := stmtAddLoc.Exec(l.URL, l.Hash)
	if err != nil {
		return nil, err
	}
	l.id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return l, nil
}
