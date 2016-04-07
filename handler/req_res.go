package handler

import (
	"encoding/json"
	"github.com/evolsnow/binding"
	"github.com/evolsnow/samaritan/common/base"
	"net/http"
)

//base response for all requests
type baseResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

func makeBaseResp(w http.ResponseWriter, r *http.Request) {
	makeResp(w, r, baseResp{Code: 0})
}

//struct to post to-do request
type postTdReq struct {
	StartTime  int64  `json:"startTime"`
	Place      string `json:"place"`
	Repeat     bool   `json:"repeat"`
	RepeatMode int    `json:"repeatMode"`
	AllDay     bool   `json:"allDay"`
	Desc       string `json:"desc"`
	Remark     string `json:"remark"`
	MissionPId string `json:"missionId"`
}

func (pt *postTdReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pt.StartTime: binding.Field{
			Form:     "startTime",
			Required: true,
		},
		&pt.Desc: binding.Field{
			Form:     "desc",
			Required: true,
		},
		&pt.Place:      "place",
		&pt.Repeat:     "repeat",
		&pt.RepeatMode: "repeatMode",
		&pt.AllDay:     "allDay",
		&pt.Remark:     "remark",
		&pt.MissionPId: "missionId",
	}
}

//struct to post to-do response
type postTdResp struct {
	baseResp
	Id string `json:"id"`
}

//new user
type postUsReq struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	StuNum     string `json:"stuNum"`
	Mail       string `json:"mail"`
	Password   string `json:"password"`
	Type       string `json:"type"`
	VerifyCode string `json:"verifyCode"`
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
		&pu.Type: binding.Field{
			Form:     "type",
			Required: true,
		},
		&pu.VerifyCode: binding.Field{
			Form:     "verifyCode",
			Required: true,
		},
		&pu.Phone:  "phone",
		&pu.Mail:   "mail",
		&pu.StuNum: "stuNum",
	}
}

type postUsResp struct {
	baseResp
	Id    string `json:"id"`
	Token string `json:"token"`
}

//new project
type postPjReq struct {
	Desc    string `json:"desc"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
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
		&pp.Private: "private",
	}
}

type postPjResp struct {
	baseResp
	Id string `json:"id"`
}

//new chat
type postPrivateChatReq struct {
	From string `json:"from"`
	To   string `json:"to"`
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

//new verify code
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

//login
type postAccessTokenReq struct {
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
	SamId    string `json:"samId"`
	Type     string `json:"type"`
	Password string `json:"password"`
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

//get method

//samId available status

type samIdStatusResp struct {
	baseResp
}

//user projects

type NestedProjects struct {
	Id          string `json:"id"` //public id
	Name        string `json:"name"`
	Desc        string `json:"desc,omitempty"` //description for the project
	CreatorId   string `json:"creatorId"`      //who created the project
	CreatorName string `json:"creatorName"`    //who created the project
	Private     bool   `json:"private"`
	Type        string `json:"type"` //joined or created
}

type userProjectsResp struct {
	baseResp
	Np []NestedProjects `json:"projects"`
}

//put method

//change password
type putPasswordReq struct {
	//Phone      string
	//SamId      string
	//Mail       string
	Password   string `json:"password"`
	Type       string `json:"type"`
	VerifyCode string `json:"verifyCode"`
}

func (pp *putPasswordReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pp.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
		&pp.Type: binding.Field{
			Form:     "type",
			Required: true,
		},
		&pp.VerifyCode: binding.Field{
			Form:     "verifyCode",
			Required: true,
		},
		//&pp.Phone: "phone",
		//&pp.Mail:  "mail",
		//&pp.SamId: "samId",
	}
}

type putPasswordResp struct {
	baseResp
}

//update to-do
type putTdReq struct {
	StartTime  int64  `json:"startTime"`
	Place      string `json:"place"`
	Repeat     bool   `json:"repeat"`
	RepeatMode int    `json:"repeatMode"`
	AllDay     bool   `json:"allDay"`
	Desc       string `json:"desc"`
	Remark     string `json:"remark"`
	MissionPId string `json:"missionId"`
	Done       bool   `json:"done"`
	FinishTime int64  `json:"finishTime"`
}

func (pt *putTdReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pt.StartTime:  "startTime",
		&pt.Desc:       "desc",
		&pt.Place:      "place",
		&pt.Repeat:     "repeat",
		&pt.RepeatMode: "repeatMode",
		&pt.AllDay:     "allDay",
		&pt.Remark:     "remark",
		&pt.MissionPId: "missionId",
		&pt.Done:       "done",
		&pt.FinishTime: "finishTime",
	}
}

type putTdResp struct {
	baseResp
}

//delete method

//delete to-do
type delTodoResp struct {
	baseResp
}

//delete project
type delProjectResp struct {
	baseResp
}

//bind json to user defined struct
func parseReq(w http.ResponseWriter, r *http.Request, ds interface{}) bool {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ds)
	if err != nil {
		base.BadReqErr(w, err.Error())
		return false
	}
	return true
}

//write user defined struct to client
func makeResp(w http.ResponseWriter, r *http.Request, src interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(src)
	if err != nil {
		base.InternalErr(w, err.Error())
	}
}
