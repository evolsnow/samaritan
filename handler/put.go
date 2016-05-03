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
	NotLoginErr         = "您还未登录"
)

// UpdatePassword updates user's password
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

// UpdateUserInfo updates other common info
func UpdateUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putUserInfoReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	user := model.InitedUser(ps.GetInt("authId"))
	if user == nil {
		base.ForbidErr(w, NotLoginErr)
		return
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Alias != "" {
		user.Alias = req.Alias
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	user.Save()
	makeResp(w, r, putUserInfoResp{})
}

// UpdateTodo updates to-do's info
func UpdateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putTdReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	tid := dbms.ReadTodoId(ps.Get("todo"))
	if tid == 0 {
		base.NotFoundErr(w, TodoNotExistErr)
		return
	}
	td := &model.Todo{
		Id:         tid,
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
		UpdateMissionStatus(w, r, ps)
	}
	if len(td.Pictures) > 0 {
		td.UpdatePics(td.Pictures)
	}
	td.Save()
	makeResp(w, r, putTdResp{})
}

// UpdateMissionStatus sets mission status to done or not
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
	m := model.InitedMission(mid)
	if m == nil {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	if !base.InIntSlice(uid, m.ReceiversId) {
		base.ForbidErr(w, NotMissionMemberErr)
		return
	}
	if req.Done && !base.InIntSlice(uid, u.GetAllCompletedMissionsId()) {
		m.CompletionNum += 100 / len(m.ReceiversId)
		u.CompleteMission(m.Id)
		if m.CompletionNum == 100 {
			m.CompletedTime = time.Now().Unix()
		}
	}
	if !req.Done && base.InIntSlice(uid, u.GetAllCompletedMissionsId()) {
		m.CompletionNum -= 100 / len(m.ReceiversId)
		u.UnCompleteMission(m.Id)
		m.CompletedTime = 0
	}
	m.ForceSave()
	makeResp(w, r, putMsStatusResp{})
}

// UpdateMissionPics updates mission's pictures
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

//AcceptMission updates user's accepted mission
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

// JoinProject updates user's joined project list
func JoinProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putJoinPjReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	pid := dbms.ReadProjectId(ps.Get("project"))
	if pid == 0 {
		base.NotFoundErr(w, ProjectNotExistErr)
		return
	}
	uid := ps.GetInt("authId")
	u := &model.User{Id: uid}
	p := &model.Project{Id: pid}
	p.AddMember(uid)
	u.JoinProject(p.Id)
	makeResp(w, r, putAcceptMsResp{})
}

// UpdateChatStatus set chat status to dealt or not
func UpdateChatStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(putCtStatusReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	uid := ps.GetInt("authId")
	u := &model.User{Id: uid}
	cid := dbms.ReadChatId(ps.Get("chat"))
	c := model.InitedChat(cid)
	if c == nil {
		base.NotFoundErr(w, ChatNotExistErr)
		return
	}
	if base.InIntSlice(c.Id, u.GetAllMsgsId()) {
		u.DealtChat(cid, req.Dealt)
	}
	makeResp(w, r, putCtStatusResp{})
}
