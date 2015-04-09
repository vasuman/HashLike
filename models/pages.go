package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/url"
	"regexp"
)

var (
	ErrNoSuchGroupID = errors.New("non-existent group ID")
	ErrInvURL        = errors.New("invalid URL")
	ErrWrongProto    = errors.New("wrong protocol")
	ErrNoDomainMatch = errors.New("URL doesn;t match any domain pattern")
	ErrNoPathMatch   = errors.New("URL doesn't match any path pattern")
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

type Pattern struct {
}

func NewPattern(s string, sep char) (*Pattern, error) {

}

func (p *Pattern) Match(t string) bool {

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
	Name         string
	Proto        protoSpec
	System       string
	SkipFragment bool
	Domains      []*regexp.Regexp
	Paths        []*regexp.Regexp
}

func (g *Group) IsValid(loc string) (string, err) {
	u, err := url.Parse(loc)
	if err != nil {
		return "", ErrInvURL
	}
	protoValid := false
	if (g.Proto&ProtoPlain) == 1 && u.Proto == "http" {
		protoValid = true
	}
	if (g.Proto&ProtoSecure) == 1 && u.Proto == "https" {
		protoValid = true
	}
	if !protoValid {
		return "", ErrWrongProto
	}
	if g.SkipFragment {
		u.Fragment = ""
	}
	return u.String(), nil
}

func AddGroup(group *Group) error {
	res, err := stmtAddGroup.Exec(group.Name, group.Proto, group.System, group.SkipFragment)
	if err != nil {
		return err
	}
	group.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	for _, domain := range group.Domains {
		_, err = stmtAddDomainPattern.Exec(group.ID, domain.String())
	}
	for _, path := range group.Paths {
		_, err = stmtAddPathPattern.Exec(group.ID, path.String())
	}
}

func getPatterns(rows *sql.Rows) (exps []*regexp.Regexp) {
	for rows.Next() {
		var str string
		_ = rows.Scan(&str)
		re := regexp.MustCompile(str)
		append(exps, re)
	}
	return
}

func GetGroup(id int64) (*Group, error) {
	g := new(Group)
	row := stmtGetGroup.QueryRow(id)
	err := row.Scan(&g.ID, &g.Name, &g.Proto, &g.System, &g.SkipFragment)
	if err == sql.ErrNoRows {
		return nil, ErrNoSuchGroupID
	}
	panicIf(err)
	dRows, err := stmtGetDomainPatterns.Query(g.ID)
	panicIf(err)
	defer dRows.Close()
	g.Domains = getPatterns(dRows)
	pRows, err := stmtGetPathPatterns.Query(g.ID)
	panicIf(err)
	defer pRows.Close()
	g.Paths = getPatterns(pRows)
	return g, nil
}
