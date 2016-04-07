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
	td := &model.Todo{Id: tid}
	td.Load()
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
	pj := &model.Project{Id: pid}
	pj.Load()
	if pj.Desc != req.Desc || pj.Private != req.Private {
		t.Error("save project error")
	}
}
