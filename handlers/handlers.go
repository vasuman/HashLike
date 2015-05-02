package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/vasuman/HashLike/res"
)

var logger *log.Logger

func getIP(r *http.Request) []byte {
	ip := net.ParseIP(r.RemoteAddr)
	//TODO: check the 'X-Real-IP' and 'X-Forwarded-For' headers.
	return ip
}

func method(meth string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != meth {
			code := http.StatusMethodNotAllowed
			txt := fmt.Sprintf("method %s not allowed here", r.Method)
			http.Error(w, txt, code)
			return
		}
		fn(w, r)
	}
}

func badRequest(w http.ResponseWriter, msg string) {
	code := http.StatusBadRequest
	http.Error(w, "bad request: "+msg, code)
}

func internalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	txt := fmt.Sprintf("internal server error - %v", err)
	http.Error(w, txt, code)
}

func renderTemplate(w http.ResponseWriter, name string, ctx interface{}) {
	var b bytes.Buffer
	err := res.Template.ExecuteTemplate(&b, name, ctx)
	if err != nil {
		logger.Printf("error rendering template:\n%v\n", err)
		internalError(w, err)
		return
	}
	io.Copy(w, &b)
}

func GetRootHandler(logInst *log.Logger) http.Handler {
	logger = logInst
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

func showURL(w http.ResponseWriter, r *http.Request) {
}
