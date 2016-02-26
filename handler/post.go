package handler

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/model"
	"github.com/mholt/binding"
	"log"
	"net/http"
)

func NewTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("calling func: NewTodo")
	req := new(postTodoRequest)
	errs := binding.Bind(r, req)
	if errs != nil {
		base.BadReqErrorHandle(w, errs.Error())
		return
	}

	//if req.OwnerId != ps.Get("userId") {
	//	base.ForbidErrorHandler(w)
	//	return
	//}
	go func() {
		td := model.Todo{
			Desc:         req.Desc,
			OwnerId:      req.OwnerId,
			Deadline:     req.Deadline,
			StartTime:    req.StartTime,
			Accomplished: req.Accomplished,
			MissionId:    req.MissionId,
		}
		td.Save()
	}()
	resp := new(postTodoResponse)
	generateResponse(w, r, resp)
}
