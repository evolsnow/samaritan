package handler

import (
	"bytes"
	"encoding/json"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/model"
	"net/http"
	"testing"
	"time"
)

func post(reqURL, auth string, src interface{}, ds interface{}) {
	var t testing.T
	body, _ := json.Marshal(src)
	//reqURL = url.QueryEscape(reqURL)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("http post err")
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(ds)
	if err != nil {
		t.Error(err)
	}
}

func TestNewUser(t *testing.T) {
	req := &postUsReq{
		Name:       "李四",
		Mail:       "abc@def.com",
		Password:   "pwd",
		Type:       "mail",
		VerifyCode: "123456",
	}
	reply := new(postUsResp)
	cache.Set("abc@def.com:code", "123456", time.Minute*5)
	post("http://127.0.0.1:8080/users", "", req, reply)
	if reply.Id == "" || reply.Token == "" {
		t.Error("create user failed")
	}

	req.VerifyCode = "654321"
	reply = new(postUsResp)
	post("http://127.0.0.1:8080/users", "", req, reply)
	if reply.Msg != CodeMismatchErr {
		t.Error("code mismatch")
	}

	req.VerifyCode = "123456"
	cache.Delete("abc@def.com:code")
	reply = new(postUsResp)
	post("http://127.0.0.1:8080/users", "", req, reply)
	if reply.Msg != ExpiredErr {
		t.Error("code expired")
	}
}

func TestNewAccessToken(t *testing.T) {
	//based on test new user
	req := &postAccessTokenReq{
		Mail:     "abc@def.com",
		Type:     "mail",
		Password: "pwd",
	}
	reply := new(postAccessTokenResp)
	post("http://127.0.0.1:8080/accessToken", "", req, reply)
	if reply.Id == "" || reply.Token == "" {
		t.Error("login failed")
	}
	req.Password = "abc"
	reply = new(postAccessTokenResp)
	post("http://127.0.0.1:8080/accessToken", "", req, reply)
	if reply.Msg != PasswordMismatchErr {
		t.Error("pwd mismatch")
	}
}

func TestNewTodo(t *testing.T) {
	req := &postTdReq{
		StartTime: time.Now().Unix(),
		Desc:      "desc",
		Place:     "one place",
		Repeat:    true,
	}
	reply := new(postTdResp)

	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)

	post("http://127.0.0.1:8080/todos", "", req, reply)
	if reply.Code != http.StatusUnauthorized {
		t.Error("unauthorized to create todo")
	}
	//normal case
	reply = new(postTdResp)
	post("http://127.0.0.1:8080/todos", auth, req, reply)
	if reply.Id == "" {
		t.Error("post todo error")
		t.FailNow()
	}
	tid := dbms.ReadTodoId(reply.Id)
	td := model.InitedTodo(tid)
	if td.Desc != req.Desc || td.Place != req.Place {
		t.Error("save todo error")
	}
}

func TestNewProject(t *testing.T) {
	req := &postPjReq{
		Name:    "name",
		Desc:    "desc",
		Private: false,
	}

	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	//normal case
	reply := new(postPjResp)
	post("http://127.0.0.1:8080/projects", auth, req, reply)
	if reply.Id == "" {
		t.Error("post projects error")
		t.FailNow()
	}
	pid := dbms.ReadProjectId(reply.Id)
	pj := model.InitedProject(pid)
	if pj.Desc != req.Desc || pj.Private != req.Private {
		t.Error("save project error")
	}
}

func TestNewInvitation(t *testing.T) {
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	req := &postProjectInvitationReq{
		Invitee:     base.HashedUserId(uid),
		ProjectId:   "e2a7af07009a48fce8b0c2646f5089d3",
		ProjectName: "pj name",
		Remark:      "remark",
	}
	auth := base.MakeToken(uid)
	reply := new(postProjectInvitationResp)
	post("http://127.0.0.1:8080/invitations/project", auth, req, reply)
	if reply.Code != 0 {
		t.Error("failed to invite to project")
	}
}

func TestNewComment(t *testing.T) {
	mPid := cache.Get("post_test_mission_pid")
	req := &postCommentReq{
		MissionPid: mPid,
	}
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	reply := new(postCommentResp)
	post("http://127.0.0.1:8080/comments", base.MakeToken(111), req, reply)
	if reply.Msg != UnableToCommentErr {
		t.Error("should be unable to comment")
	}
	post("http://127.0.0.1:8080/comments", auth, req, reply)
	if reply.Code != 0 {
		t.Error("failed to comment")
	}
	m := &model.Mission{Id: dbms.ReadMissionId(mPid)}
	if len(m.GetComments()) == 0 {
		t.Error("failed to save comment")
	}
}

func TestNewMission(t *testing.T) {
	req := &postMissionReq{
		Name:        "ms name",
		Desc:        "ms desc",
		ReceiversId: []string{base.HashedUserId(1), base.HashedUserId(2)},
		ProjectId:   cache.Get("post_test_project_pid"),
	}
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	reply := new(postMissionResp)
	post("http://127.0.0.1:8080/missions", auth, req, reply)
	if reply.Code != 0 {
		t.Error("create mission error:", reply.Msg)
	}
}

func TestNewMissionInvitation(t *testing.T) {
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	req := &postMissionInvitationReq{
		Invitee:     base.HashedUserId(uid),
		MissionId:   "e2a7af07009a48fce8b0c2646f5089d3",
		MissionName: "ms name",
		Remark:      "remark",
	}
	auth := base.MakeToken(uid)
	reply := new(postMissionInvitationResp)
	post("http://127.0.0.1:8080/invitations/mission", auth, req, reply)
	if reply.Code != 0 {
		t.Error("failed to invite to mission")
	}
}
