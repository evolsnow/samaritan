package handler

import (
	"github.com/evolsnow/binding"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/model"
	"net/http"
	"time"
)

const (
	NotMissionMemberErr = "您还未接受此任务"
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
	makeResp(w, r, putPasswordResp{})
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
		Pictures:   req.Pictures,
	}
	if td.GetOwner().Id != ps.GetInt("authId") {
		base.ForbidErr(w, BelongErr)
		return
	}
	if td.Done {
		td.Finish()
	}
	if len(td.Pictures) > 0 {
		td.UpdatePics(td.Pictures)
	}
	td.Save()
	makeResp(w, r, putTdResp{})
}

//
//func UpdateTodoPics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	req := new(putTdPicReq)
//	errs := binding.Bind(r, req)
//	if errs.Handle(w) {
//		return
//	}
//	log.DebugJson(req)
//	tid := dbms.ReadTodoId(ps.Get("todo"))
//	t := &model.Todo{Id:tid}
//	t.UpdatePics(req.Pictures)
//	makeResp(w, r, putTdPicResp{})
//}

func UpdateMissionStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putMsStatusReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	uid := ps.GetInt("authId")
	u := &model.User{Id: uid}
	mid := dbms.ReadMissionId(ps.Get("mission"))
	if mid == 0 {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	m := &model.Mission{Id: mid}
	receivers := m.GetReceiversId()
	if !base.InIntSlice(uid, receivers) {
		base.ForbidErr(w, NotMissionMemberErr)
		return
	}
	m.Load()
	if req.Done && !base.InIntSlice(uid, u.GetAllCompletedMissionsId()) {
		m.CompletionNum += 100 / len(receivers)
		u.CompleteMission(m.Id)
		if m.CompletionNum == 100 {
			m.CompletedTime = time.Now().Unix()
		}
	}
	if !req.Done && base.InIntSlice(uid, u.GetAllCompletedMissionsId()) {
		m.CompletionNum -= 100 / len(receivers)
		u.UnCompleteMission(m.Id)
		m.CompletedTime = 0
	}
	m.ForceSave()
	makeResp(w, r, putMsStatusResp{})
}

func UpdateMissionPics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putMsPicReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	mid := dbms.ReadMissionId(ps.Get("mission"))
	m := &model.Mission{Id: mid}
	m.UpdatePics(req.Pictures)
	makeResp(w, r, putMsPicResp{})
}

func AcceptMission(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putAcceptMsReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	mid := dbms.ReadMissionId(ps.Get("mission"))
	if mid == 0 {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	uid := ps.GetInt("authId")
	u := &model.User{Id: uid}
	m := &model.Mission{Id: mid}
	m.AddReceiver(uid)
	u.AcceptMission(m.Id)
	makeResp(w, r, putAcceptMsResp{})
}
