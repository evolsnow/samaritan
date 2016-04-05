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

func post(reqURL, auth string, src []byte, ds interface{}) {
	var t testing.T
	//reqURL = url.QueryEscape(reqURL)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", reqURL, bytes.NewReader(src))
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

func TestNewTodo(t *testing.T) {
	req := &postTdReq{
		StartTime: time.Now().Unix(),
		Desc:      "desc",
		Place:     "one place",
		Repeat:    true,
	}
	reply := new(postTdResp)

	src, _ := json.Marshal(req)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)

	post("http://127.0.0.1:8080/todos", "", src, reply)
	if reply.Code != http.StatusUnauthorized {
		t.Error("unauthorized to create todo")
	}
	//normal case
	reply = new(postTdResp)
	post("http://127.0.0.1:8080/todos", auth, src, reply)
	t.Log(reply)
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
