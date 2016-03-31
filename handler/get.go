package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"net/http"
)

const (
	LengthErr  = "length should be 4-8"
	CharsetErr = "charset should be a-z, A-Z, 0-9 or _"
	ExistErr   = "already registered"
)

func SamIdStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	samId := ps.Get("samId")
	resp := new(samIdStatusResp)
	if len(samId) > 8 || len(samId) < 4 {
		resp.Msg = LengthErr
		log.DebugJson(resp)
		makeResp(w, r, resp)
		return
	}
	if !base.ValidSamId(samId) {
		resp.Msg = CharsetErr
		log.DebugJson(resp)
		makeResp(w, r, resp)
		return
	}
	if dbms.ReadIfSamIdExist(samId) {
		resp.Available = false
		resp.Msg = ExistErr
		log.DebugJson(resp)
		makeResp(w, r, resp)
		return
	}
	resp.Available = true
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
