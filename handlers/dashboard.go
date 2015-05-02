package handlers

import (
	"net/http"

	"github.com/vasuman/HashLike/db"
	"github.com/vasuman/HashLike/pow"
)

func getSystems() (ret [][]string) {
	for k, v := range pow.Systems {
		ret = append(ret, []string{k, v.Desc()})
	}
	return
}

type HeaderParams struct {
}

func showDashboard(w http.ResponseWriter, r *http.Request) {
	type dashParams struct {
		Header  struct{}
		Groups  []*db.Group
		Systems [][]string
	}
	params := new(dashParams)
	var err error
	params.Groups, err = db.ListGroups()
	params.Systems = getSystems()
	if err != nil {
		internalError(w, err)
		return
	}
	renderTemplate(w, "dashboard", params)
}

func showGroup(w http.ResponseWriter, r *http.Request) {

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
	if form.Get("skip-fragment") != "" {
		g.SkipFragment = true
	}
	sys := form.Get("sys")
	if _, ok := pow.Systems[sys]; !ok {
		badRequest(w, "invalid proof-of-work system identifier")
		return
	}
	g.System = sys
	logger.Printf("built group %+v\n", g)
	err = db.AddGroup(g)
	if err != nil {
		internalError(w, err)
		return
	}
	http.Redirect(w, r, "show?key="+g.Key, http.StatusSeeOther)
}

func editGroup(w http.ResponseWriter, r *http.Request) {
}
