package handler

import (
	"encoding/json"
	"github.com/evolsnow/binding"
	"github.com/evolsnow/samaritan/common/base"
	"net/http"
)

//base response for all requests
type baseResp struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func makeBaseResp(w http.ResponseWriter, r *http.Request) {
	makeResp(w, r, baseResp{Code: http.StatusOK})
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
	Name       string
	Phone      string
	Mail       string
	Password   string
	Source     string
	VerifyCode string
}

func (pu *postUsReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pu.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
		&pu.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
		&pu.Source: binding.Field{
			Form:     "source",
			Required: true,
		},
		&pu.VerifyCode: binding.Field{
			Form:     "verifyCode",
			Required: true,
		},
		&pu.Phone: "phone",
		&pu.Mail:  "mail",
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

type postVerifyCodeReq struct {
	To  string `json:"to"`
	Use string `json:"use"`
}

func (pvc *postVerifyCodeReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pvc.To: binding.Field{
			Form:     "to",
			Required: true,
		},
		&pvc.Use: binding.Field{
			Form:     "use",
			Required: true,
		},
	}
}

type postVerifyCodeResp struct {
	baseResp
}

type postAccessTokenReq struct {
	Phone    string
	Mail     string
	SamId    string
	Type     string
	Password string
}

func (pat *postAccessTokenReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pat.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
		&pat.Type: binding.Field{
			Form:     "type",
			Required: true,
		},
		&pat.Phone: "phone",
		&pat.Mail:  "mail",
		&pat.SamId: "samId",
	}
}

type postAccessTokenResp struct {
	baseResp
	Id    string `json:"id"`
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
