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
	BelongErr = "请检查所登录的账户"
)

func DeleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tid := dbms.ReadTodoId(ps.Get("todo"))
	uid := ps.GetInt("authId")
	log.Debug("authId:", uid)
	td := &model.Todo{
		Id: tid,
	}
	if td.GetOwner().Id != uid {
		base.ForbidErr(w, BelongErr)
		return
	}
	td.Remove()
	makeBaseResp(w, r)
}

func DeleteProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := dbms.ReadProjectId(ps.Get("project"))
	uid := ps.GetInt("authId")
	p := &model.Project{
		Id:        pid,
		CreatorId: uid,
	}
	if p.GetCreator().Id != uid {
		base.ForbidErr(w, BelongErr)
		return
	}
	p.Remove()
	makeBaseResp(w, r)
}
