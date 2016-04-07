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
	nps := make([]NestedProjects, len(pjs))
	var createdOrJoined string
	for i, p := range pjs {
		log.DebugJson(p)
		if p.CreatorId == uid {
			createdOrJoined = "created"
		} else {
			createdOrJoined = "joined"
		}
		np := NestedProjects{
			Id:          p.Pid,
			Name:        p.Name,
			Desc:        p.Desc,
			CreatorId:   base.HashedUserId(p.CreatorId),
			CreatorName: p.GetCreator().Name,
			Private:     p.Private,
			Type:        createdOrJoined,
		}
		nps[i] = np
	}
	resp.Np = nps
	log.DebugJson(resp)
	makeResp(w, r, resp)
}
