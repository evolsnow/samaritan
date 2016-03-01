package handler

import (
	"fmt"
	"github.com/evolsnow/binding"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/model"
	"net/http"
)

func NewUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postUsReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.Debug(req)
	us := model.User{
		Phone:    req.Phone,
		Password: base.EncryptedPassword(req.Password),
		Name:     req.Name,
	}
	go us.CreateAvatar()
	//assign id to user
	us.Save()
	//return jwt token
	token := base.NewToken(us.Id)
	resp := new(postUsResp)
	resp.Id = base.HashedUserId(us.Id)
	resp.Token = token
	go createToken(us.Id, token)
	log.Debug(resp)
	makeResp(w, r, resp)
}

func NewTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postTdReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	uid := ps.GetInt("userId")
	td := model.Todo{
		OwnerId:      uid,
		Desc:         req.Desc,
		StartTime:    req.StartTime,
		TaskTime:     req.TaskTime,
		Place:        req.Place,
		Repeat:       req.Repeat,
		RepeatPeriod: req.RepeatPeriod,
		MissionId:    req.ProjectId,
	}
	td.Save()
	resp := new(postTdResp)
	resp.Id = td.Pid
	makeResp(w, r, resp)
}

func NewProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postPjReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	pj := model.Project{
		Desc: req.Desc,
		Name: req.Name,
	}
	pj.CreatorId = ps.GetInt("userId")
	//get pid
	pj.Save()
	resp := new(postPjResp)
	resp.Id = pj.Pid
	makeResp(w, r, resp)
}

func NewPrivateChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postPrivateChatReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	fd := ps.GetInt("userId")
	td := model.ReadUserId(req.To)
	raw := ""
	if fd < td {
		raw = fmt.Sprintf("%d&%d", fd, td)

	} else {
		raw = fmt.Sprintf("%d&%d", td, fd)
	}
	resp := new(postPrivateChatResp)
	resp.PrivateChatId = base.NewPrivateChatId(raw)
	//go createPrivateChatRecord(resp.PrivateChatId, fd, td)
	makeResp(w, r, resp)
}
