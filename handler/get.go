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

	ProjectNotExistErr = "项目不存在"
	MissionNotExistErr = "任务不存在"

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
			Id:          p.Pid,
			Name:        p.Name,
			Desc:        p.Desc,
			CreatorId:   base.HashedUserId(p.CreatorId),
			CreatorName: p.GetCreator().Name,
			Private:     p.Private,
			Type:        createdOrJoined,
			Members:     p.GetMembersName(),
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
	u := &model.User{Id: uid}
	u.Load()
	resp := userSearchResp{
		Name:   u.Name,
		Id:     u.Pid,
		Avatar: u.FullAvatarUrl(),
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func ProjectMissionList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := dbms.ReadProjectId(ps.Get("project"))
	if pid == 0 {
		base.NotFoundErr(w, ProjectNotExistErr)
		return
	}
	resp := new(projectMissionsResp)
	p := &model.Project{Id: pid}
	ms := p.GetMissions()
	uid := ps.GetInt("authId")
	u := &model.User{Id: uid}
	u.Load()
	userAcceptedMissions := u.GetAllAcceptedMissionsId()
	if !base.InIntSlice(uid, p.GetMembersId()) {
		base.ForbidErr(w, NotProjectMemberErr)
		return
	}
	nms := make([]NestedMission, len(ms))
	for i, v := range ms {
		if !base.InIntSlice(v.Id, userAcceptedMissions) {
			continue
		}
		nm := NestedMission{
			Id:            v.Pid,
			Name:          v.Name,
			Desc:          v.Desc,
			Deadline:      v.Deadline,
			Pictures:      v.GetPictures(),
			CreatorName:   u.Name,
			CreatorId:     u.Pid,
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
	if mid == 0 {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	resp := new(missionCommentResp)
	m := &model.Mission{Id: mid}
	cms := m.GetComments()

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
	mid := dbms.ReadProjectId(ps.Get("mission"))
	if mid == 0 {
		base.NotFoundErr(w, MissionNotExistErr)
		return
	}
	m := &model.Mission{Id: mid}
	m.Load()
	receivers := m.GetReceiversId()
	receiversPid := make([]string, len(receivers))
	for i, v := range receivers {
		receiversPid[i] = base.HashedUserId(v)
	}
	resp := &missionDetailResp{
		Id:            ps.Get("mission"),
		CreateTime:    m.CreateTime,
		Name:          m.Name,
		Desc:          m.Desc,
		Deadline:      m.Deadline,
		Pictures:      m.GetPictures(),
		PublisherId:   base.HashedUserId(m.PublisherId),
		ReceiversId:   receiversPid,
		CompletionNum: m.CompletionNum,
		CompletedTime: m.CompletedTime,
		ProjectId:     base.HashedProjectId(m.ProjectId),
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}

func MakeUploadToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := &QiNiuUpTokenResp{
		Token:  base.QiNiuUploadToken(ps.Get("key")),
		Expire: base.QiNiuExpire,
	}
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
