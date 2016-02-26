package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/base"
	"github.com/mholt/binding"
	"net/http"
)

//base response for all requests
type baseResp struct {
	Code  int   `json:"code,omitempty"`
	Error error `json:"error,omitempty"`
}

func makeBaseResp(w http.ResponseWriter, r *http.Request) {
	makeResp(w, r, baseResp{Code: 200})
}

//struct to post to-do request
type postTdReq struct {
	StartTime uint64
	Deadline  uint64
	Desc      string
	Done      bool
	MissionId int
}

func (pt *postTdReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pt.Desc: binding.Field{
			Form:     "desc",
			Required: true,
		},
		&pt.Deadline: binding.Field{
			Form:     "deadline",
			Required: true,
		},
		&pt.StartTime: "startTime",
		&pt.Done:      "done",
		&pt.MissionId: "missionId",
	}
}

//struct to post to-do response
type postTdResp struct {
	baseResp
}

type postUsReq struct {
	Phone    string
	Password string
}

func (pu *postUsReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pu.Phone: binding.Field{
			Form:     "phone",
			Required: true,
		},
		&pu.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
		//&pu.StartTime:    "startTime",
		//&pu.Done: "done",
		//&pu.MissionId:    "missionId",
	}
}

type postUsResp struct {
	baseResp
	Token string `json:"token"`
}

//bind json to user defined struct
func parseReq(w http.ResponseWriter, r *http.Request, ds interface{}) bool {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ds)
	if err != nil {
		base.SetError(w, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

//write user defined struct to client
func makeResp(w http.ResponseWriter, r *http.Request, src interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(src)
	if err != nil {
		base.SetError(w, err.Error(), http.StatusInternalServerError)
	}
}
