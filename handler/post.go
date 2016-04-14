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
	UnknownCodeTypeErr = "未知验证码类型"
	ExpiredErr         = "验证码未发送或已过期"
	CodeMismatchErr    = "验证码不匹配"

	UnknownUseErr    = "未知验证码用途"
	UnknownSourceErr = "未知的发送渠道"
	InvalidPhoneErr  = "非法的手机号格式"
	InvalidMailErr   = "非法的邮箱地址格式"

	NotRegisteredErr     = "用户未注册"
	AlreadyRegisteredErr = "用户已注册"
	PasswordMismatchErr  = "密码不匹配"

	UserNotExistErr = "用户不存在"

	UnableToCommentErr = "不是该任务发布者或者接收者,无权评论"
)

const (
	InvitedToJoinProject = "%s 邀请你加入项目: %s"
	InvitedToJoinMission = "%s 邀请你接受任务: %s"
	DeliverMission       = "%s 发布了一个任务: %s"
)

func NewDeviceToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postDtReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	dbms.CreateDeviceIndex(ps.GetInt("authId"), req.DeviceToken)
	makeResp(w, r, postDtResp{})
}

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
		base.BadReqErr(w, UnknownCodeTypeErr)
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
	//assign id to user
	us.Save()
	go us.CreateAvatar()
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
	uid := ps.GetInt("authId")
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

//func NewTodoPics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	req := new(postPicReq)
//	errs := binding.Bind(r, req)
//	if errs.Handle(w) {
//		return
//	}
//	log.DebugJson(req)
//	uid := ps.GetInt("authId")
//	tid := dbms.ReadTodoId(ps.Get("todo"))
//
//}

func NewProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postPjReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	pj := &model.Project{
		CreatorId: ps.GetInt("authId"),
		Desc:      req.Desc,
		Name:      req.Name,
		Private:   req.Private,
	}
	pj.CreatorId = ps.GetInt("authId")
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
	fd := ps.GetInt("authId")
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
		if dbms.ReadUserIdWithIndex(req.To, source) != 0 {
			base.ForbidErr(w, AlreadyRegisteredErr)
			return
		}
		title = "帐号注册"
		body = "您的注册验证码是： " + code + "，有效期为5分钟。"
		text = "【GoDo日程】" + body
	case "resetPasswd":
		if dbms.ReadUserIdWithIndex(req.To, source) == 0 {
			base.NotFoundErr(w, NotRegisteredErr)
			return
		}
		title = "重置密码"
		body = "您正在申请重置密码，验证码为： " + code + "，有效期为5分钟。（如非本人操作，请尽快查看账户操作情况）"
		text = "【GoDo日程】" + body

	default:
		base.BadReqErr(w, UnknownUseErr)
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
		base.BadReqErr(w, UnknownSourceErr)
		return
	}
	go cache.Set(req.To+":code", code, time.Minute*5)
	makeResp(w, r, postVerifyCodeResp{})
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

func NewProjectInvitation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postProjectInvitationReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	uid := ps.GetInt("authId")
	user := &model.User{Id: uid}
	user.Load()
	go func() {
		msg := fmt.Sprintf(InvitedToJoinProject, user.Name, req.ProjectName)
		payload := make(map[string]string)
		payload["invitor"] = user.Pid
		payload["projectId"] = req.ProjectId
		payload["remark"] = req.Remark
		push := &model.Chat{
			Type:      model.InvitedToProject,
			Target:    req.Invitee,
			Msg:       msg,
			ExtraInfo: payload,
		}
		log.DebugJson(push)
		push.Response()
	}()
	makeResp(w, r, postProjectInvitationResp{})
}

func NewMission(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postMissionReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	uid := ps.GetInt("authId")
	user := &model.User{Id: uid}
	user.Load()
	m := &model.Mission{
		PublisherId: uid,
		Name:        req.Name,
		Desc:        req.Desc,
		ProjectId:   dbms.ReadProjectId(req.ProjectId),
	}
	m.Save()
	go func() {
		msg := fmt.Sprintf(DeliverMission, user.Name, m.Name)
		payload := make(map[string]string)
		payload["invitor"] = user.Pid
		payload["missionId"] = m.Pid
		push := &model.Chat{
			Type:      model.InvitedToMission,
			To:        req.ReceiversId,
			Msg:       msg,
			ExtraInfo: payload,
		}
		log.DebugJson(push)
		push.Response()
	}()
	resp := &postMissionResp{
		Id: m.Pid,
	}
	makeResp(w, r, resp)
}

func NewMissionInvitation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postMissionInvitationReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	uid := ps.GetInt("authId")
	user := &model.User{Id: uid}
	user.Load()
	go func() {
		msg := fmt.Sprintf(InvitedToJoinMission, user.Name, req.MissionName)
		payload := make(map[string]string)
		payload["invitor"] = user.Pid
		payload["missionId"] = req.MissionId
		payload["remark"] = req.Remark
		push := &model.Chat{
			Type:      model.InvitedToMission,
			To:        []string{req.Invitee},
			Msg:       msg,
			ExtraInfo: payload,
		}
		log.DebugJson(push)
		push.Response()
	}()
	makeResp(w, r, postMissionInvitationResp{})
}

func NewComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(postCommentReq)
	errs := binding.Bind(r, req)
	if errs.Handle(w) {
		return
	}
	log.DebugJson(req)
	uid := ps.GetInt("authId")
	user := &model.User{Id: uid}
	user.Load()
	mid := dbms.ReadMissionId(req.MissionPid)
	m := &model.Mission{Id: mid}
	m.Load()
	if uid != m.PublisherId || !base.InIntSlice(uid, m.ReceiversId) {
		log.Debug(uid, m.PublisherId, m.ReceiversId)
		base.ForbidErr(w, UnableToCommentErr)
		return
	}
	cm := &model.Comment{
		CriticPid:  user.Pid,
		CriticName: user.Name,
		MissionPid: req.MissionPid,
	}
	cm.Save()
	resp := new(postCommentResp)
	resp.Id = cm.Pid
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
