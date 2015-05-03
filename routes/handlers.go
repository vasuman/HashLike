package routes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/vasuman/HashLike/res"
)

const stylePrefix = "/styles/"

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
	logger.Printf("error: %v\n", err)
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
	d.HandleFunc("/group/", method("GET", listGroups))
	d.HandleFunc("/group/show", method("GET", showGroup))
	d.HandleFunc("/group/add", method("POST", addGroup))
	d.HandleFunc("/group/patterns", method("POST", setPatterns))
	d.HandleFunc("/group/check", method("POST", checkURL))
	//TODO: dashboard authentication
	r.Handle("/group/", d)
	a := http.NewServeMux()
	a.HandleFunc("/api/count", method("GET", getCount))
	a.HandleFunc("/api/like", method("POST", newChallenge))
	a.HandleFunc("/api/verify", method("POST", verifySoln))
	r.Handle("/api/", a)
	//TODO: url public/private
	r.HandleFunc("/url", method("GET", showURL))
	r.HandleFunc(stylePrefix, method("GET", serveStyle))
	return r
}

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "#L")
}

func serveStyle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, stylePrefix)
	s, ok := res.Styles[path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Content-Type", "text/css")
	w.Write(s)
}

func showURL(w http.ResponseWriter, r *http.Request) {
}
