package handler

import (
	"fmt"
	"github.com/evolsnow/binding"
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/common/rpc"
	"github.com/evolsnow/samaritan/model"
	"net/http"
	"time"
)

func NewUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postUsReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	var code, info, source string
	if req.Source == "phone" {
		code = cache.Get(req.Phone + ":code")
		info, source = req.Phone, "phone"
	} else if req.Source == "mail" {
		code = cache.Get(req.Mail + ":code")
		info, source = req.Mail, "mail"
	} else {
		base.ForbidErrorHandler(w, "unknown register source")
		return
	}
	if code == "" {
		base.ForbidErrorHandler(w, "code has expired")
		return
	}
	if code != req.VerifyCode {
		base.ForbidErrorHandler(w, "code mismatch")
		return
	}
	us := model.User{
		Phone:    req.Phone,
		Email:    req.Mail,
		Password: base.EncryptedPassword(req.Password),
		Name:     req.Name,
	}
	go us.CreateAvatar()
	//assign id to user
	us.Save()
	//create user login index
	go dbms.CreateLoginIndex(us.Id, info, source)
	//return jwt token
	token := base.MakeToken(us.Id)
	resp := new(postUsResp)
	resp.Id = base.HashedUserId(us.Id)
	resp.Token = token
	go dbms.CreateToken(us.Id, token)
	log.DebugJson(resp)
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
		OwnerId:   uid,
		Desc:      req.Desc,
		StartTime: req.StartTime,
		Place:     req.Place,
		Repeat:    req.Repeat,
		MissionId: req.ProjectId,
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

func NewVerifyCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postVerifyCodeReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	source := ps.Get("source")
	code := base.RandomCodeSix()
	var title, body, text string
	switch req.Use {
	case "register":
		title = "帐号注册"
		body = "您的注册验证码是： " + code + "，有效期为5分钟。"
		text = "【GoDo日程】" + body
	case "forgetPasswd":
		title = "找回密码"
		body = "您正在申请找回密码，验证码为： " + code + "，有效期为5分钟。（如非本人操作，请尽快查看账户操作情况）"
		text = "【GoDo日程】" + body
	case "resetPasswd":
		title = "重置密码"
		body = "您正在申请重置密码，验证码为： " + code + "，有效期为5分钟。（如非本人操作，请尽快查看账户操作情况）"
		text = "【GoDo日程】" + body
	default:
		base.BadReqErrHandle(w, "unknown use")
		return
	}
	if source == "sms" {
		if !base.ValidPhone(req.To) {
			base.BadReqErrHandle(w, "invalid phone number")
			return
		}
		go rpc.SendSMS(req.To, text)
	} else if source == "mail" {
		if !base.ValidMail(req.To) {
			base.BadReqErrHandle(w, "invalid mail address")
			return
		}
		go rpc.SendMail(req.To, title, body)
	} else {
		base.BadReqErrHandle(w, "unknown source")
		return
	}
	go cache.Set(req.To+":code", code, time.Minute*5)
	makeBaseResp(w, r)
}

func NewAccessToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postAccessTokenReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	var uid int
	switch req.Type {
	case "phone":
		uid = dbms.ReadLoginUid(req.Phone, req.Type)
	case "mail":
		uid = dbms.ReadLoginUid(req.Mail, req.Type)
	case "samId":
		uid = dbms.ReadLoginUid(req.SamId, req.Type)
	}
	if uid == 0 {
		base.SetError(w, "user not registered", http.StatusNotFound)
		return
	}
	us := new(model.User)
	us.Id = uid
	if base.EncryptedPassword(req.Password) != us.GetPassword() {
		base.ForbidErrorHandler(w, "password mismatch")
		return
	}
	resp := new(postAccessTokenResp)
	resp.Id = base.HashedUserId(us.Id)
	resp.Token = dbms.ReadToken(us.Id)
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
