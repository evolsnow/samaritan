package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/model"
	"github.com/mholt/binding"
	"net/http"
)

// base response for all requests
type baseResp struct {
	Code  int   `json:"code"`
	Error error `json:"error, omitempty"`
}

type postTodoRequest struct {
	StartTime    uint64 `json:"startTime, omitempty"`
	Deadline     uint64 `json:"deadline"`
	Desc         string `json:"desc"`
	OwnerId      int    `json:"ownerId"`
	Accomplished bool   `json:"accomplished, omitempty"`
	MissionId    int    `json:"missionId, omitempty"`
}

func (pt *postTodoRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pt.Desc: binding.Field{
			Form:     "desc",
			Required: true,
		},
		&pt.OwnerId: binding.Field{
			Form:     "ownerId",
			Required: true,
		},
		&pt.Deadline: binding.Field{
			Form:     "deadline",
			Required: true,
		},
		&pt.StartTime:    "startTime",
		&pt.Accomplished: "accomplished",
		&pt.MissionId:    "missionId",
	}
}

func (pt *postTodoRequest) SaveTodo() (err error) {
	td := model.Todo{
		Desc:         pt.Desc,
		OwnerId:      pt.OwnerId,
		Deadline:     pt.Deadline,
		StartTime:    pt.StartTime,
		Accomplished: pt.Accomplished,
		MissionId:    pt.MissionId,
	}
	return td.Save()
}

type postTodoResponse struct {
	baseResp
}

//bind json to user defined struct
func parseRequest(w http.ResponseWriter, r *http.Request, ds interface{}) bool {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ds)
	if err != nil {
		base.SetError(w, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

//write user defined struct to client
func generateResponse(w http.ResponseWriter, r *http.Request, src interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(src)
	if err != nil {
		base.SetError(w, err.Error(), http.StatusInternalServerError)
	}
}
