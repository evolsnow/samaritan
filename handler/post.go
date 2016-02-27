package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/model"
	"github.com/mholt/binding"
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
	makeResp(w, r, resp)
}

func NewTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postTdReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	uid, _ := strconv.Atoi(ps.Get("userId"))
	go func() {
		td := model.Todo{
			OwnerId:   uid,
			Desc:      req.Desc,
			Deadline:  req.Deadline,
			StartTime: req.StartTime,
			Done:      req.Done,
			MissionId: req.MissionId,
		}
		td.Save()
	}()
	makeBaseResp(w, r)
}
