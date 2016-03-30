package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"net/http"
)

func SamIdStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	samId := ps.Get("samId")
	resp := new(samIdStatusResp)
	if len(samId) > 8 || len(samId) < 4 {
		resp.Msg = "length should be 4-8"
		makeResp(w, r, resp)
		return
	}
	if !base.ValidSamId(samId) {
		resp.Msg = "charset should be a-z, A-Z, 0-9 or _"
		makeResp(w, r, resp)
		return
	}
	if dbms.ReadIfSamIdExist(samId) {
		resp.Available = false
		resp.Msg = "already registered"
		makeResp(w, r, resp)
		return
	}
	resp.Available = true
	makeResp(w, r, resp)
}
