package db

import (
	"errors"
	"strings"
)

const (
	Star      = "*"
	PathSep   = "/"
	DomainSep = "."
)

type PathPattern struct {
	Parts    []*PartMatcher
	Complete bool
}

func (p *PathPattern) Matches(path string) bool {
	parts := strings.Split(path, PathSep)
	if len(parts) < len(p.Parts) {
		return false
	}
	if p.Complete && len(parts) != len(p.Parts) {
		return false
	}
	for i, pm := range p.Parts {
		if !pm.Match(parts[i]) {
			return false
		}
	}
	return true
}

func ParsePath(s string) (*PathPattern, error) {
	p := new(PathPattern)
	if strings.HasSuffix(s, "$") {
		p.Complete = true
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "/") {
		//trailing slash
		s += "*"
	}
	pts := strings.Split(s, PathSep)
	for _, pt := range pts {
		part, err := parsePart(pt)
		if err != nil {
			return nil, err
		}
		p.Parts = append(p.Parts, part)
	}
	return p, nil
}

type DomainPattern struct {
	Parts []*PartMatcher
}

func (d *DomainPattern) Matches(domain string) bool {
	parts := strings.Split(domain, DomainSep)
	if len(parts) != len(d.Parts) {
		return false
	}
	for i, pm := range d.Parts {
		if !pm.Match(parts[i]) {
			return false
		}
	}
	return true
}

func ParseDomain(s string) (*DomainPattern, error) {
	d := new(DomainPattern)
	pts := strings.Split(s, DomainSep)
	for _, pt := range pts {
		part, err := parsePart(pt)
		if err != nil {
			return nil, err
		}
		d.Parts = append(d.Parts, part)
	}
	return d, nil
}

const (
	allMatch byte = iota
	prefixMatch
	suffixMatch
	substrMatch
	emptyMatch
	exactMatch
)

type PartMatcher struct {
	Kind  byte
	Param string
}

func (p *PartMatcher) Match(target string) bool {
	switch p.Kind {
	case allMatch:
		return true
	case emptyMatch:
		return len(target) == 0
	case prefixMatch:
		return strings.HasPrefix(target, p.Param)
	case suffixMatch:
		return strings.HasSuffix(target, p.Param)
	case substrMatch:
		return strings.Contains(target, p.Param)
	case exactMatch:
		return target == p.Param
	}
	return false
}

func parsePart(p string) (pm *PartMatcher, err error) {
	pm = new(PartMatcher)
	if len(p) == 0 {
		pm.Kind = emptyMatch
		return
	}
	if p == Star {
		pm.Kind = allMatch
		return
	}
	first := strings.HasPrefix(p, Star)
	last := strings.HasSuffix(p, Star)
	switch count := strings.Count(p, Star); count {
	case 0:
		pm.Kind = exactMatch
		pm.Param = p
		return
	case 1:
		if first {
			pm.Kind = suffixMatch
			pm.Param = p[1:]
			return
		} else if last {
			pm.Kind = prefixMatch
			pm.Param = p[:len(p)-1]
			return
		}
		err = errors.New("not prefix or suffix")
	case 2:
		if first && last {
			pm.Kind = substrMatch
			pm.Param = p[1 : len(p)-1]
			return
		}
		err = errors.New("not substr")
	default:
		err = errors.New("too many stars")
	}
	return
}
