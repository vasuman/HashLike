package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func decodeReq(r reqResp, v interface{}) bool {
	dec := json.NewDecoder(r.r.Body)
	err := dec.Decode(v)
	if err != nil {
		errStr := fmt.Sprintf("error decoding request - %v", err)
		internalError(r.w, errStr)
		return false
	}
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

	ok := decodeReq(reqResp{w, r}, req)
	if !ok {
		return
	}
	b, err := json.Marshal(resp)
	w.Write(b)
}

func verifySoln(w http.ResponseWriter, r *http.Request) {

}
