package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/vasuman/HashLike/db"
	"github.com/vasuman/HashLike/pow"
)

type headerParams struct {
	Title  string
	Styles []string
}

type reqResp struct {
	w http.ResponseWriter
	r *http.Request
}

func groupWrap(ga func(reqResp, *db.Group, url.Values)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		form := r.PostForm
		key := form.Get("key")
		if key == "" {
			badRequest(w, "need a valid 'key'")
			return
		}
		g, err := db.GetGroup(key)
		if err != nil {
			internalError(w, err)
			return
		}
		ga(reqResp{w, r}, g, form)
	}
}

func listGroups(w http.ResponseWriter, r *http.Request) {
	type dashParams struct {
		Header  headerParams
		Groups  []*db.Group
		Systems map[string]string
	}
	groups, err := db.ListGroups()
	if err != nil {
		internalError(w, err)
		return
	}
	params := &dashParams{
		headerParams{
			"Groups",
			[]string{"group"},
		},
		groups,
		pow.Desc,
	}
	renderTemplate(w, "groupList", params)
}

func showGroup(w http.ResponseWriter, r *http.Request) {
	type showGroupParams struct {
		Header headerParams
		Group  *db.Group
	}
	q := r.URL.Query()
	key := q.Get("key")
	if key == "" {
		badRequest(w, "need a valid 'key'")
		return
	}
	params := &showGroupParams{
		headerParams{
			"Group",
			[]string{"group"},
		},
		nil,
	}
	var err error
	params.Group, err = db.GetGroup(key)
	if err != nil {
		internalError(w, err)
		return
	}
	renderTemplate(w, "groupShow", params)
}

func addGroup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	form := r.PostForm
	g := new(db.Group)
	name := form.Get("name")
	if name == "" {
		badRequest(w, "need a valid 'name'")
		return
	}
	g.Name = name
	proto, err := db.ProtoFromString(form.Get("proto"))
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	g.Proto = proto
	if form.Get("strip-fragment") != "" {
		g.StripFragment = true
	}
	sys := form.Get("sys")
	if _, ok := pow.Systems[sys]; !ok {
		badRequest(w, "invalid proof-of-work system identifier")
		return
	}
	g.System = sys
	err = db.AddGroup(g)
	if err != nil {
		internalError(w, err)
		return
	}
	logger.Printf("built group %+v\n", g)
	http.Redirect(w, r, "show?key="+g.Key, http.StatusSeeOther)
}

func setPatterns(p reqResp, g *db.Group, form url.Values) {
	const newline = "\n"
	var (
		dm  *db.DomainMatcher
		pm  *db.PathMatcher
		err error
	)
	domains := strings.Split(form.Get("domains"), newline)
	paths := strings.Split(form.Get("paths"), newline)
	g.Domains, g.Paths = nil, nil
	for _, dom := range domains {
		dom = strings.TrimSpace(dom)
		if dom == "" {
			continue
		}
		dm, err = db.ParseDomain(dom)
		if err != nil {
			break
		}
		g.Domains = append(g.Domains, dm)
	}
	if err != nil {
		badRequest(p.w, "invalid domain pattern")
		return
	}
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		pm, err = db.ParsePath(path)
		if err != nil {
			break
		}
		g.Paths = append(g.Paths, pm)
	}
	if err != nil {
		badRequest(p.w, "invalid path pattern")
		return
	}
	err = db.UpdateGroup(g)
	if err != nil {
		internalError(p.w, err)
		return
	}
	fmt.Fprintf(p.w, "updated patterns")
}

func checkURL(p reqResp, g *db.Group, form url.Values) {
	//	http.Redirect(w, r, "show?key="+g.Key, http.StatusSeeOther)
	_, err := g.IsValid(form.Get("url"))
	if err != nil {
		fmt.Fprintf(p.w, "error - %v", err)
		return
	}
	fmt.Fprintf(p.w, "yes")
}

func deleteGroup(p reqResp, g *db.Group, form url.Values) {
	if g.Name != form.Get("name") {
		badRequest(p.w, "group name mismatch")
		return
	}
	db.DeleteGroup(g)
	http.Redirect(p.w, p.r, "", http.StatusSeeOther)
}
