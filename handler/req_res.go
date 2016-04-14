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

////upload pic
//type postPicReq struct {
//	PicUrl []string `json:"picUrl"`
//}
//
//type postPicResp struct {
//	baseResp
//}

//upload device token
type postDtReq struct {
	DeviceToken string `json:"deviceToken"`
}

func (pd *postDtReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pd.DeviceToken: binding.Field{
			Form:     "deviceToken",
			Required: true,
		},
	}
}

type postDtResp struct {
	baseResp
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
		&pp.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
		&pp.Desc:    "desc",
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

//project invitation
type postProjectInvitationReq struct {
	Invitee     string `json:"invitee"`
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	Remark      string `json:"remark"`
}

func (pi *postProjectInvitationReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pi.Invitee: binding.Field{
			Form:     "invitee",
			Required: true,
		},
		&pi.ProjectId: binding.Field{
			Form:     "projectId",
			Required: true,
		},
		&pi.ProjectName: binding.Field{
			Form:     "projectName",
			Required: true,
		},
		&pi.Remark: "remark",
	}
}

type postProjectInvitationResp struct {
	baseResp
}

//mission invitation
type postMissionInvitationReq struct {
	Invitee     string `json:"invitee"`
	MissionId   string `json:"missionId"`
	MissionName string `json:"missionName"`
	Remark      string `json:"remark"`
}

func (pi *postMissionInvitationReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pi.Invitee: binding.Field{
			Form:     "invitee",
			Required: true,
		},
		&pi.MissionId: binding.Field{
			Form:     "missionId",
			Required: true,
		},
		&pi.MissionName: binding.Field{
			Form:     "missionName",
			Required: true,
		},
		&pi.Remark: "remark",
	}
}

type postMissionInvitationResp struct {
	baseResp
}

//mission
type postMissionReq struct {
	Name        string   `json:"name,omitempty"`
	Desc        string   `json:"desc,omitempty"`
	ReceiversId []string `json:"receiversId,omitempty"`
	ProjectId   string   `json:"projectId,omitempty"`
}

func (pm *postMissionReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pm.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
		&pm.ProjectId:   "projectId",
		&pm.Desc:        "desc",
		&pm.ReceiversId: "receiversId",
	}
}

type postMissionResp struct {
	baseResp
	Id string `json:"id"`
}

//mission comment
type postCommentReq struct {
	MissionPid string `json:"mission"`
}

func (pc *postCommentReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&pc.MissionPid: binding.Field{
			Form:     "mission",
			Required: true,
		},
	}
}

type postCommentResp struct {
	baseResp
	Id string `json:"id"`
}

//get method

//samId available status

type samIdStatusResp struct {
	baseResp
}

//user projects

type NestedProject struct {
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
	Np []NestedProject `json:"projects"`
}

//search user
type userSearchResp struct {
	baseResp
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

//project missions

type NestedMission struct {
	Id            string `json:"id"` //public id
	Name          string `json:"name"`
	Desc          string `json:"desc,omitempty"` //description for the project
	CreatorName   string `json:"creatorName,omitempty"`
	CreatorId     string `json:"creatorId,omitempty"`
	CreateTime    int64  `json:"createTime,omitempty"`
	completionNum int    `json:"completionNum"`
}

type projectMissionsResp struct {
	baseResp
	Nm []NestedMission `json:"missions"`
}

//mission comments

type NestedComment struct {
	Id         string `json:"id"`
	CreateTime int64  `json:"createTime"`
	CriticPid  string `json:"userId"`
	CriticName string `json:"userName"`
}

type missionCommentResp struct {
	baseResp
	Nm []NestedComment `json:"comments"`
}

//mission detail
type missionDetailResp struct {
	baseResp
	Id            string   `json:"id"`
	CreateTime    int64    `json:"createTime,omitempty"`
	Name          string   `json:"name,omitempty"`
	Desc          string   `json:"desc,omitempty"`
	PublisherId   string   `json:"publisherId,omitempty"`
	ReceiversId   []string `json:"receiversId,omitempty"`
	CompletionNum int      `json:"completionNum,omitempty"`
	CompletedTime int64    `json:"completedTime,omitempty"`
	ProjectId     string   `json:"projectId,omitempty"`
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

//update mission
type putMsStatusReq struct {
	Done bool `json:"done"`
}

func (pm *putMsStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		//&pm.Done: binding.Field{
		//	Form:     "done",
		//	Required: true,
		//},
		&pm.Done: "done",
	}
}

type putMsStatusResp struct {
	baseResp
}

//accept mission
type putAcceptMsReq struct {
}

func (pm *putAcceptMsReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
	//&pm.Done: "done",
	}
}

type putAcceptMsResp struct {
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
