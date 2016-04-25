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
	BelongErr        = "请检查所登录的账户"
	MissionRemainErr = "项目包含任务，暂无法删除"
)

func DeleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tid := dbms.ReadTodoId(ps.Get("todo"))
	uid := ps.GetInt("authId")
	log.Debug("authId:", uid)
	td := model.InitedTodo(tid)
	if td.OwnerId != uid {
		base.ForbidErr(w, BelongErr)
		return
	}
	td.Remove()
	makeBaseResp(w, r)
}

func DeleteMission(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mid := dbms.ReadMissionId(ps.Get("mission"))
	uid := ps.GetInt("authId")
	log.Debug("authId:", uid)
	m := model.InitedMission(mid)
	m.Sync()
	if m.PublisherId != uid {
		base.ForbidErr(w, BelongErr)
		return
	}
	m.Remove()
	makeBaseResp(w, r)
}

func DeleteProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := dbms.ReadProjectId(ps.Get("project"))
	uid := ps.GetInt("authId")
	p := model.InitedProject(pid)
	if p.CreatorId != uid {
		base.ForbidErr(w, BelongErr)
		return
	}
	if len(p.GetMissions()) > 0 {
		base.ForbidErr(w, MissionRemainErr)
		return
	}
	p.Remove()
	makeBaseResp(w, r)
}
