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

const (
	UnknownTypeErr  = "unknown type"
	ExpiredErr      = "code has expired"
	CodeMismatchErr = "code mismatch"

	UnknownUseErr    = "unknown use"
	UnknownSourceErr = "unknown source"
	InvalidPhoneErr  = "invalid phone number"
	InvalidMailErr   = "invalid mail address"

	NotRegisteredErr    = "user not registered"
	PasswordMismatchErr = "password mismatch"
)

func NewUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postUsReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	var code, info, source string
	if req.Type == "phone" {
		code = cache.GetSet(req.Phone+":code", "")
		info, source = req.Phone, "phone"
	} else if req.Type == "mail" {
		//code = cache.GetSet(req.Mail+":code", "")
		code = cache.Get(req.Mail + ":code")
		info, source = req.Mail, "mail"
	} else {
		base.BadReqErr(w, UnknownTypeErr)
		return
	}
	if code == "" {
		base.ForbidErr(w, ExpiredErr)
		return
	}
	if code != req.VerifyCode {
		base.ForbidErr(w, CodeMismatchErr)
		return
	}
	us := &model.User{
		Phone:      req.Phone,
		Email:      req.Mail,
		Password:   base.EncryptedPassword(req.Password),
		Name:       req.Name,
		StudentNum: req.StuNum,
	}
	go us.CreateAvatar()
	//assign id to user
	us.Save()
	//create user login/search index
	go dbms.CreateSearchIndex(us.Id, info, source)
	//return jwt token and public user id
	resp := new(postUsResp)
	resp.Id = us.Pid
	resp.Token = base.MakeToken(us.Id)
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
	td := &model.Todo{
		OwnerId:    uid,
		StartTime:  req.StartTime,
		Place:      req.Place,
		Repeat:     req.Repeat,
		RepeatMode: req.RepeatMode,
		AllDay:     req.AllDay,
		Desc:       req.Desc,
		Remark:     req.Remark,
		MissionId:  dbms.ReadMissionId(req.MissionPId),
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
	pj := &model.Project{
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
	td := dbms.ReadUserId(req.To)
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
		base.BadReqErr(w, UnknownSourceErr)
		return
	}
	if source == "sms" {
		if !base.ValidPhone(req.To) {
			base.BadReqErr(w, InvalidPhoneErr)
			return
		}
		go rpc.SendSMS(req.To, text)
	} else if source == "mail" {
		if !base.ValidMail(req.To) {
			base.BadReqErr(w, InvalidMailErr)
			return
		}
		go rpc.SendMail(req.To, title, body)
	} else {
		base.BadReqErr(w, UnknownUseErr)
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
		uid = dbms.ReadUserIdWithIndex(req.Phone, req.Type)
	case "mail":
		uid = dbms.ReadUserIdWithIndex(req.Mail, req.Type)
	case "samId":
		uid = dbms.ReadUserIdWithIndex(req.SamId, req.Type)
	}
	if uid == 0 {
		base.NotFoundErr(w, NotRegisteredErr)
		return
	}
	us := new(model.User)
	us.Id = uid
	if base.EncryptedPassword(req.Password) != us.GetPassword() {
		base.ForbidErr(w, PasswordMismatchErr)
		return
	}
	resp := new(postAccessTokenResp)
	resp.Id = base.HashedUserId(us.Id)
	resp.Token = base.MakeToken(us.Id)
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
