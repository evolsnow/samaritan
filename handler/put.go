package handler

import (
	"github.com/evolsnow/binding"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/model"
	"net/http"
)

func UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putPasswordReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	identity := ps.Get("identity")
	var code string
	if req.Type == "phone" {
		code = cache.GetSet(req.Phone+":code", "")
	} else if req.Type == "mail" {
		code = cache.GetSet(req.Mail+":code", "")
	} else {
		base.BadReqErr(w, "unknown type")
		return
	}
	if code == "" {
		base.ForbidErr(w, "code has expired")
		return
	}
	if code != req.VerifyCode {
		base.ForbidErr(w, "code mismatch")
		return
	}
	uid := dbms.ReadUidIndex(identity, req.Type)
	if uid == 0 {
		base.NotFoundErr(w, "user not found")
		return
	}
	us := new(model.User)
	us.Id = uid
	us.Password = base.EncryptedPassword(req.Password)
	go us.Save()
	makeBaseResp(w, r)
}
