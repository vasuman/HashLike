package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"math/rand"
	"net/url"
	"regexp"
)

var (
	ErrNoSuchGroupKey = errors.New("non-existent group key")
	ErrInvURL         = errors.New("invalid URL")
	ErrInvProto       = errors.New("wrong protocol")
	ErrNoDomainMatch  = errors.New("URL doesn't match any domain pattern")
	ErrNoPathMatch    = errors.New("URL doesn't match any path pattern")
)

var (
	stmtAddDomainPattern  *sql.Stmt
	stmtGetDomainPatterns *sql.Stmt
	stmtAddPathPattern    *sql.Stmt
	stmtGetPathPatterns   *sql.Stmt
	stmtAddGroup          *sql.Stmt
	stmtGetGroup          *sql.Stmt
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type protoSpec int

func (p *protoSpec) Scan(value interface{}) error {
	*p = protoSpec(value.(int))
	return nil
}

func (p protoSpec) Value() (driver.Value, error) {
	return int(p), nil
}

const (
	ProtoPlain protoSpec = iota + 1
	ProtoSecure
	ProtoBoth
)

type Group struct {
	ID           int64
	Key          string
	Name         string
	Proto        protoSpec
	System       string
	SkipFragment bool
	Domains      []*regexp.Regexp
	Paths        []*regexp.Regexp
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
	//TODO: check domain and path patterns
	return u.String(), nil
}

const (
	Chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	GroupKeyLen = 6
)

func genGroupKey() string {
	id := make([]byte, GroupKeyLen)
	for i := range id {
		r := rand.Intn(len(Chars))
		id[i] = Chars[r]
	}
	return string(id)
}

func keyUnique(key string) bool {
	_, err := stmtGetGroup.Query(key)
	return err == sql.ErrNoRows
}

func AddGroup(group *Group) error {
	key := genGroupKey()
	for !keyUnique(key) {
		key = genGroupKey()
	}
	res, err := stmtAddGroup.Exec(key, group.Name, group.Proto, group.System, group.SkipFragment)
	if err != nil {
		return err
	}
	group.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	for _, domain := range group.Domains {
		_, err = stmtAddDomainPattern.Exec(group.ID, domain.String())
		if err != nil {
			return err
		}
	}
	for _, path := range group.Paths {
		_, err = stmtAddPathPattern.Exec(group.ID, path.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func getPatterns(rows *sql.Rows) (exps []*regexp.Regexp) {
	for rows.Next() {
		var str string
		_ = rows.Scan(&str)
		re := regexp.MustCompile(str)
		exps = append(exps, re)
	}
	return
}

func GetGroup(key string) (*Group, error) {
	g := new(Group)
	row := stmtGetGroup.QueryRow(key)
	err := row.Scan(&g.ID, &g.Key, &g.Name, &g.Proto, &g.System, &g.SkipFragment)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoSuchGroupKey
		}
		return nil, err
	}
	dRows, err := stmtGetDomainPatterns.Query(g.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer dRows.Close()
	g.Domains = getPatterns(dRows)
	pRows, err := stmtGetPathPatterns.Query(g.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer pRows.Close()
	g.Paths = getPatterns(pRows)
	return g, nil
}
