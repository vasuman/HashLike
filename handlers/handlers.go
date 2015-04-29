package handlers

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/vasuman/HashLike/res"
)

var challengeTimeout time.Duration

func getIP(r *http.Request) []byte {
	ip := net.ParseIP(r.RemoteAddr)
	//TODO: Check the 'X-Real-IP' and 'X-Forwarded-For' headers.
	return ip
}

func method(meth string, fn http.HandlerFunc) http.HandlerFunc {
	txt := fmt.Sprintf("method %s not allowed here", meth)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != meth {
			code := http.StatusMethodNotAllowed
			http.Error(w, txt, code)
			return
		}
		fn(w, r)
	}
}

func internalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	txt := fmt.Sprintf("internal server error - %v", err)
	http.Error(w, txt, code)
}

func renderTmpl(w http.ResponseWriter, name string, ctx interface{}) {
	err := res.Template.ExecuteTemplate(w, name, ctx)
	if err != nil {
		internalError(w, err)
	}
}

func GetRootHandler() http.Handler {
	res.Setup()
	r := http.NewServeMux()
	d := http.NewServeMux()
	r.HandleFunc("/", root)
	d.HandleFunc("/dash/", method("GET", showDashboard))
	d.HandleFunc("/dash/group/show", method("GET", showGroup))
	d.HandleFunc("/dash/group/add", method("POST", addGroup))
	d.HandleFunc("/dash/group/edit", method("POST", editGroup))
	//TODO: dashboard authentication
	r.Handle("/dash/", d)
	a := http.NewServeMux()
	a.HandleFunc("/api/count", method("GET", getCount))
	a.HandleFunc("/api/like", method("POST", newChallenge))
	a.HandleFunc("/api/verify", method("POST", verifySoln))
	r.Handle("/api/", a)
	//TODO: url public/private
	r.HandleFunc("/url", method("GET", showURL))
	return r
}

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "#L")
}

func showDashboard(w http.ResponseWriter, r *http.Request) {
	renderTmpl(w, "dashboard", nil)
}

func showGroup(w http.ResponseWriter, r *http.Request) {
}

func addGroup(w http.ResponseWriter, r *http.Request) {
}

func editGroup(w http.ResponseWriter, r *http.Request) {
}

func showURL(w http.ResponseWriter, r *http.Request) {
}

func getCount(w http.ResponseWriter, r *http.Request) {
}

func newChallenge(w http.ResponseWriter, r *http.Request) {
}

func verifySoln(w http.ResponseWriter, r *http.Request) {
}
