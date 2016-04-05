package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"net/http"
)

const (
	LengthErr  = "长度应为4-8位"
	CharsetErr = "仅支持a-z, A-Z, 0-9 以及 _"
	ExistErr   = "已经被注册"
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

//func UserProjectList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	uid:=ps.GetInt("authId")
//	pjType := r.URL.Query().Get("type")
//	us:=&model.User{Id:uid}
//	us.
//}
