package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"github.com/evolsnow/samaritan/model"
	"net/http"
)

const (
	LengthErr  = "长度应为4-8位"
	CharsetErr = "仅支持a-z, A-Z, 0-9 以及 _"
	ExistErr   = "已经被注册"

	UnknownTypeErr = "未知类型"

	ProjectNotExistErr  = "项目不存在"
	MissionNotExistErr  = "任务不存在"
	TodoNotExistErr     = "Todo不存在"
	ChatNotExistErr     = "消息不存在"
	NotProjectMemberErr = "不是本项目成员,无法查看项目"
)

func SamIdStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	samId := ps.Get("samId")
	resp := new(samIdStatusResp)
	if len(samId) > 8 || len(samId) < 4 {
		resp.Code = 1
		resp.Msg = LengthErr
		log.DebugJson(resp)
		makeResp(w, r, resp)
		return
	}
	if !base.ValidSamId(samId) {
		resp.Code = 2
		resp.Msg = CharsetErr
		log.DebugJson(resp)
		makeResp(w, r, resp)
		return
	}
	if dbms.ReadIfSamIdExist(samId) {
		resp.Code = 3
		resp.Msg = ExistErr
		log.DebugJson(resp)
		makeResp(w, r, resp)
		return
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func UserProjectList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.GetInt("authId")
	pjType := r.URL.Query().Get("type")
	us := &model.User{Id: uid}
	resp := new(userProjectsResp)
	var pjs []model.Project
	switch pjType {
	case "joined":
		pjs = us.GetJoinedProjects()
	case "created":
		pjs = us.GetCreatedProjects()
	case "":
		pjs = us.GetAllProjects()
	default:
		base.BadReqErr(w, UnknownTypeErr)
		return
	}
	nps := make([]NestedProject, len(pjs))
	var createdOrJoined string
	for i, p := range pjs {
		if p.CreatorId == uid {
			createdOrJoined = "created"
		} else {
			createdOrJoined = "joined"
		}
		np := NestedProject{
			Id:           p.Pid,
			Name:         p.Name,
			Desc:         p.Desc,
			CreatorId:    base.HashedUserId(p.CreatorId),
			CreatorName:  p.GetCreator().Name,
			Private:      p.Private,
			Type:         createdOrJoined,
			MembersCount: len(p.GetMembersName()),
		}
		nps[i] = np
	}
	resp.Np = nps
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func SearchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userMail := ps.Get("user")
	uid := dbms.ReadUserIdWithIndex(userMail, "mail")
	if uid == 0 {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	u := model.InitedUser(uid)
	resp := userSearchResp{
		Name:   u.Name,
		Id:     u.Pid,
		Avatar: u.Avatar,
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func ProjectDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := dbms.ReadProjectId(ps.Get("project"))
	p := model.InitedProject(pid)
	if p == nil {
		base.NotFoundErr(w, ProjectNotExistErr)
		return
	}
	uid := ps.GetInt("authId")
	u := model.InitedUser(uid)
	if u == nil {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	members := p.GetMembers()
	allMembers := make([]userModel, len(members))
	for i, mem := range members {
		allMembers[i] = userModel{
			Id:     base.HashedUserId(mem.Id),
			Name:   mem.Name,
			Avatar: mem.FullAvatarUrl(),
		}
	}
	resp := &projectDetailResp{
		Id:          ps.Get("project"),
		CreateTime:  p.CreateTime,
		Name:        p.Name,
		Desc:        p.Desc,
		CreatorId:   base.HashedUserId(p.CreatorId),
		CreatorName: p.GetCreator().Name,
		Private:     p.Private,
		Members:     allMembers,
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func ProjectMissionList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := dbms.ReadProjectId(ps.Get("project"))
	resp := new(projectMissionsResp)
	p := model.InitedProject(pid)
	if p == nil {
		base.NotFoundErr(w, ProjectNotExistErr)
		return
	}
	owner := model.InitedUser(p.CreatorId)
	if owner == nil {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	ms := p.GetMissions()
	uid := ps.GetInt("authId")
	u := model.InitedUser(uid)
	if u == nil {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	userAcceptedMissions := u.GetAllAcceptedMissionsId()
	if !base.InIntSlice(uid, p.MembersId) {
		base.ForbidErr(w, NotProjectMemberErr)
		return
	}
	nms := make([]NestedMission, len(ms))
	for i, v := range ms {
		//if !base.InIntSlice(v.Id, userAcceptedMissions) {
		//	continue
		//}
		v.Sync()
		nm := NestedMission{
			Id:            v.Pid,
			Name:          v.Name,
			Desc:          v.Desc,
			Deadline:      v.Deadline,
			Pictures:      v.Pictures,
			Accepted:      base.InIntSlice(v.Id, userAcceptedMissions),
			ReceiversName: v.GetReceiversName(),
			CreatorName:   owner.Name,
			CreatorId:     owner.Pid,
			CreatorAvatar: owner.Avatar,
			CreateTime:    v.CreateTime,
			CompletionNum: v.CompletionNum,
		}
		nms[i] = nm
	}
	resp.Nm = nms
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func MissionCommentList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mid := dbms.ReadProjectId(ps.Get("mission"))
	m := model.InitedMission(mid)
	if m == nil {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	cms := m.Comments
	resp := new(missionCommentResp)
	ncs := make([]NestedComment, len(cms))
	for i, v := range cms {
		nc := NestedComment{
			Id:         v.Pid,
			CreateTime: v.CreateTime,
			CriticPid:  v.CriticPid,
			CriticName: v.CriticName,
		}
		ncs[i] = nc
	}
	resp.Nm = ncs
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func MissionDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mid := dbms.ReadMissionId(ps.Get("mission"))
	m := model.InitedMission(mid)
	if m == nil {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	uid := ps.GetInt("authId")
	u := model.InitedUser(uid)
	if u == nil {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	userAcceptedMissions := u.GetAllAcceptedMissionsId()
	resp := &missionDetailResp{
		Id:            ps.Get("mission"),
		CreateTime:    m.CreateTime,
		Name:          m.Name,
		Desc:          m.Desc,
		Deadline:      m.Deadline,
		Accepted:      base.InIntSlice(m.Id, userAcceptedMissions),
		Pictures:      m.Pictures,
		PublisherId:   base.HashedUserId(m.PublisherId),
		ReceiversName: m.GetReceiversName(),
		CompletionNum: m.CompletionNum,
		CompletedTime: m.CompletedTime,
		ProjectId:     base.HashedProjectId(m.ProjectId),
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func MakeUploadToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := &QiNiuUpTokenResp{
		Token:  base.QiNiuUploadToken(),
		Expire: base.QiNiuExpire,
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func OfflineMsgs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.GetInt("authId")
	u := model.InitedUser(uid)
	if u == nil {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	resp := new(getOfflineMsgResp)
	msgs := u.GetAllOfflineMsg()
	nms := make([]NestedMsg, len(msgs))
	for i, v := range msgs {
		nc := NestedMsg{
			Id:        v.Pid,
			Time:      v.Timestamp,
			Msg:       v.Msg,
			Dealt:     v.Dealt,
			Type:      v.Type,
			ExtraInfo: v.ExtraInfo,
		}
		nms[i] = nc
	}
	resp.Msgs = nms
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.GetInt("authId")
	user := model.InitedUser(uid)
	if user == nil {
		base.NotFoundErr(w, UserNotExistErr)
		return
	}
	resp := &personalInfoResp{
		Id:     user.Pid,
		Avatar: user.Avatar,
		Name:   user.Name,
		Alias:  user.Alias,
		Mail:   user.Email,
		StuNum: user.StudentNum,
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
