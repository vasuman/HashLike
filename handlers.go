package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/vasuman/HashLike/models"
	"github.com/vasuman/HashLike/pow"
)

var challengeTimeout time.Duration

func handler() http.Handler {
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/count", getCountHandler)
	rootMux.HandleFunc("/challenge", challengeHandler)
	rootMux.HandleFunc("/solution", solutionHandler)
	return rootMux
}

func getIP(r *http.Request) []byte {
	//TODO: Check the 'X-Real-IP' and 'X-Forwarded-For' headers.
	ip := net.ParseIP(r.RemoteAddr)
	return ip
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
	sys, ok := pow.Systems[cReq.System]
	if !ok {
		badRequest(w, "invalid system")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	remoteAddr := getIP(r)
	if remoteAddr == nil {
		logger.Printf("ERR: request %+v has no IP\n")
	}
	challenge := sys.Challenge(cReq.URL, remoteAddr)
	expiry := time.Now().Add(CacheExpiryDelay)
	cResp := &response{
		Challenge: encode(challenge),
		Expiry:    time.Now().Add(time.Minute * 2),
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(cResp)
	if err != nil {
		logger.Printf("ERR: encoding response - %v\n", err)
	}
}

func solutionHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Challenge string `json:"challenge"`
		Nonce     int    `json:"nonce"`
	}

	type response struct {
		Ok    bool    `json:"ok"`
		Value float64 `json:"value"`
		Error string  `json:"error,omitempty"`
	}

	if r.Method != "POST" {
		wrongMethod(w)
		return
	}

}
