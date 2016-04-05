package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/model"
	"net/http"
)

const (
	BelongErr = "the todo doesn't belong to this user"
)

func DeleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tid := dbms.ReadTodoId(ps.Get("todo"))
	uid := ps.GetInt("authId")
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
