package handler

import (
	"encoding/json"
	"github.com/evolsnow/binding"
	"github.com/evolsnow/samaritan/base"
	"github.com/evolsnow/samaritan/common/caches"
	"net/http"
)

//get cache
var cache = caches.NewCache()

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
	StartTime int64
	Desc      string
	Repeat    bool
	Place     string
	ProjectId int
}

func (pt *postTdReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pt.Desc: binding.Field{
			Form:     "desc",
			Required: true,
		},
		&pt.StartTime: binding.Field{
			Form:     "startTime",
			Required: true,
		},
		&pt.Place:     "place",
		&pt.Repeat:    "repeat",
		&pt.ProjectId: "projectId",
	}
}

//struct to post to-do response
type postTdResp struct {
	baseResp
	Id string `json:"id"`
}

type postUsReq struct {
	Name     string
	Phone    string
	Password string
}

func (pu *postUsReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pu.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
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
	Id    string `json:"id"`
	Token string `json:"token"`
}

type postPjReq struct {
	Desc string
	Name string
}

func (pp *postPjReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pp.Desc: binding.Field{
			Form:     "desc",
			Required: true,
		},
		&pp.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
	}
}

type postPjResp struct {
	baseResp
	Id string `json:"id"`
}

type postPrivateChatReq struct {
	From string
	To   string
}

func (ppc *postPrivateChatReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&ppc.From: binding.Field{
			Form:     "from",
			Required: true,
		},
		&ppc.To: binding.Field{
			Form:     "to",
			Required: true,
		},
	}
}

type postPrivateChatResp struct {
	baseResp
	PrivateChatId string `json:"chatId"`
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
