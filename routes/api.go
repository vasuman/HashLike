package routes

import (
	"encoding/json"
	"net/http"
)

func decodeReq(r reqResp, v interface{}) bool {
	dec := json.NewDecoder(r.r.Body)
	err := dec.Decode(v)
	if err != nil {
		internalError(r.w, err)
		return false
	}
	return true
}

func sendResp(r reqResp, v interface{}) bool {
	return true
}

func getCount(w http.ResponseWriter, r *http.Request) {
}

func newChallenge(w http.ResponseWriter, r *http.Request) {
	type request struct {
		URL string `json:"url"`
	}
	type response struct {
		Err error `json:"err"`
	}
	var (
		req  = new(request)
		resp = new(response)
	)
	ok := decodeReq(reqResp{w, r}, req)
	if !ok {
		return
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(b)
}

func verifySoln(w http.ResponseWriter, r *http.Request) {

}
