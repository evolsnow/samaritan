package handler

import (
	"fmt"
	"github.com/evolsnow/binding"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/model"
	"net/http"
	"strconv"
)

func NewUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postUsReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	us := model.User{
		Phone:    req.Phone,
		Password: base.HashedPassword(req.Password),
	}
	go us.CreateAvatar()
	//assign id to user
	us.Save()
	//return jwt token
	resp := new(postUsResp)
	resp.Token = base.NewToken(us.Id)
	resp.Id = base.HashedUserId(us.Id)
	makeResp(w, r, resp)
}

func NewTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postTdReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	uid, _ := strconv.Atoi(ps.Get("userId"))
	td := model.Todo{
		OwnerId:   uid,
		Desc:      req.Desc,
		Deadline:  req.Deadline,
		StartTime: req.StartTime,
		Done:      req.Done,
		ProjectId: req.MissionId,
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
	pj.PublisherId, _ = strconv.Atoi(ps.Get("userId"))
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
	fd, _ := strconv.Atoi(ps.Get("userId"))
	td, _ := readUserId(req.To)
	raw := ""
	if fd < td {
		raw = fmt.Sprintf("%d&%d", fd, td)

	} else {
		raw = fmt.Sprintf("%d&%d", td, fd)
	}
	resp := new(postPrivateChatResp)
	resp.PrivateChatId = base.NewPrivateChatId(raw)
	go createPrivateConvRecord(resp.PrivateChatId, fd, td)
	makeResp(w, r, resp)
}
