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

func put(reqURL, auth string, src interface{}, ds interface{}) {
	var t testing.T
	body, _ := json.Marshal(src)
	//reqURL = url.QueryEscape(reqURL)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", reqURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("http put err")
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(ds)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdatePassword(t *testing.T) {
	req := &putPasswordReq{
		Password:   "newpwd",
		Type:       "mail",
		VerifyCode: "123456",
	}
	reply := new(putPasswordResp)
	put("http://127.0.0.1:8080/users/password/gsc1215225@gmail.com", "", req, reply)
	if reply.Code != 200 {
		t.Error("update failed")
	}
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	u := &model.User{
		Id: uid,
	}
	//wait redis
	//time.Sleep(time.Second)
	if u.GetPassword() != base.EncryptedPassword(req.Password) {
		//t.Error(u.GetPassword())
		//t.Error(base.EncryptedPassword(req.Password))
		t.Error("pwd not change")
	}

	req.VerifyCode = "000000"
	cache.Set("gsc1215225@gmail.com:code", "123456", time.Minute*5)
	put("http://127.0.0.1:8080/users/password/gsc1215225@gmail.com", "", req, reply)
	if reply.Msg != CodeMismatchErr {
		t.Error("code mismatch")
	}

	req.VerifyCode = "123456"
	cache.Set("gsc@gmail.com:code", "123456", time.Minute*5)
	put("http://127.0.0.1:8080/users/password/gsc@gmail.com", "", req, reply)
	if reply.Msg != NotRegisteredErr {
		t.Error("not registed")
	}
}

func TestUpdateTodo(t *testing.T) {
	req := &putTdReq{
		Place:  "new place",
		Repeat: true,
	}
	reply := new(putTdResp)
	tPid := cache.Get("put_test_todo_pid")
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)

	//unauthorized
	put("http://127.0.0.1:8080/todos/"+tPid, "", req, reply)
	if reply.Code != http.StatusUnauthorized {
		t.Error("unauthorized to update this todo")
	}
	//belong err
	put("http://127.0.0.1:8080/todos/"+tPid, base.MakeToken(111), req, reply)
	if reply.Code != http.StatusForbidden || reply.Msg != BelongErr {
		t.Error("forbidden to update other user's todo:", reply.Msg)
	}
	//normal case
	put("http://127.0.0.1:8080/todos/"+tPid, auth, req, reply)
	if reply.Code != 200 {
		t.Error("update todo failed")
	}
	td := model.Todo{Id: dbms.ReadTodoId(tPid)}
	td.Load()
	if td.Place != req.Place {
		t.Error("todo place not changed")
	}
}
