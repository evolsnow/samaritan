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
	code := cache.GetSet(identity+":code", "")
	if code == "" {
		base.ForbidErr(w, ExpiredErr)
		return
	}
	if code != req.VerifyCode {
		base.ForbidErr(w, CodeMismatchErr)
		return
	}
	uid := dbms.ReadUserIdWithIndex(identity, req.Type)
	if uid == 0 {
		base.NotFoundErr(w, NotRegisteredErr)
		return
	}
	us := &model.User{
		Id:       uid,
		Password: base.EncryptedPassword(req.Password),
	}
	us.Save()
	makeBaseResp(w, r)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putTdReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	td := &model.Todo{
		Id:         dbms.ReadTodoId(ps.Get("todo")),
		StartTime:  req.StartTime,
		Place:      req.Place,
		Repeat:     req.Repeat,
		RepeatMode: req.RepeatMode,
		AllDay:     req.AllDay,
		Desc:       req.Desc,
		Remark:     req.Remark,
		MissionId:  dbms.ReadMissionId(req.MissionPId),
		Done:       req.Done,
		FinishTime: req.FinishTime,
	}
	if td.GetOwner().Id != ps.GetInt("authId") {
		base.ForbidErr(w, BelongErr)
		return
	}
	if td.Done {
		td.Finish()
	}
	td.Save()
	makeBaseResp(w, r)
}
