package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vasuman/HashLike/models"
	"github.com/vasuman/HashLike/pow"
)

func getRootHandler(staticDir string) http.Handler {
	const sPrefix = "/static/"
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/api/count", getCountHandler)
	rootMux.HandleFunc("/api/challenge", challengeHandler)
	rootMux.HandleFunc("/api/solution", solutionHandler)
	staticServer := http.StripPrefix(sPrefix, http.FileServer(http.Dir(staticDir)))
	rootMux.Handle(sPrefix, staticServer)
	return rootMux
}

func badRequest(w http.ResponseWriter, msg string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "bad request: %s", msg)
}

func wrongMethod(w http.ResponseWriter) {
	code := http.StatusMethodNotAllowed
	txt := "method not allowed"
	http.Error(w, txt, code)
}

func internalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	txt := fmt.Sprintf("internal server error - %v", err)
	http.Error(w, txt, code)
}

func getSystem(ident string) pow.POW {
	const hashcash = "HC"
	switch ident {
	case hashcash:
		return new(pow.Hashcash)
	}
	return nil
}

func getCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		wrongMethod(w)
		return
	}
	query := r.URL.Query()
	url := query.Get("url")
	if url == "" {
		badRequest(w, "need a url")
		return
	}
	sIdent := query.Get("sys")
	if sIdent == "" {
		badRequest(w, "need a sys")
		return
	}
	loc, err := models.GetLocation(url)
	if err != nil {
		badRequest(w, "url doesn't exist")
		return
	}
	fmt.Fprintf(w, "%x", loc.Hash)
}

func challengeHandler(w http.ResponseWriter, r *http.Request) {
	encode := func(b []byte) string {
		return base64.StdEncoding.EncodeToString(b)
	}

	type request struct {
		URL    string `json:"url"`
		System string `json:"system"`
	}

	type response struct {
		Expiry time.Time `json:"expires"`
		// Base64 encoded byte slice
		Challenge string `json:"challenge"`
	}

	if r.Method != "POST" {
		wrongMethod(w)
		return
	}
	dec := json.NewDecoder(r.Body)
	cReq := new(request)
	err := dec.Decode(cReq)
	if err != nil {
		badRequest(w, "bad json")
		return
	}
	logger.Printf("got %+v from %s\n", cReq, r.RemoteAddr)
	sys := getSystem(cReq.System)
	if sys == nil {
		badRequest(w, "invalid system")
		return
	}
	w.WriteHeader(http.StatusOK)
	challenge := sys.Challenge(r)
	cResp := &response{
		Challenge: encode(challenge),
		Expiry:    time.Now().Add(time.Minute * 2),
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(cResp)
	if err != nil {
		logger.Printf("error encoding response - %v", err)
	}
}

func solutionHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
	}
	if r.Method != "POST" {
		wrongMethod(w)
		return
	}
}
