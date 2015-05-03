package routes

import (
	"net/http"

	"github.com/vasuman/HashLike/db"
	"github.com/vasuman/HashLike/pow"
)

type headerParams struct {
	Title  string
	Styles []string
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
	logger.Printf("built group %+v\n", g)
	err = db.AddGroup(g)
	if err != nil {
		internalError(w, err)
		return
	}
	http.Redirect(w, r, "show?key="+g.Key, http.StatusSeeOther)
}

func setPatterns(w http.ResponseWriter, r *http.Request) {
}

func checkURL(w http.ResponseWriter, r *http.Request) {
}
