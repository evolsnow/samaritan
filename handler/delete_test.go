package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"net/http"
	"testing"
)

func del(reqURL string, auth string, ds interface{}) {
	var t testing.T
	//reqURL = url.QueryEscape(reqURL)
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", reqURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("http delete err")
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(ds)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteTodo(t *testing.T) {
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	tPid := cache.Get("delete_test_todo_pid")
	reply := new(delTodoResp)
	//unauthorized
	del("http://127.0.0.1:8080/todos/"+tPid, "", reply)
	if reply.Code != http.StatusUnauthorized {
		t.Error("unauthorized to delete this todo")
	}
	//belong err
	del("http://127.0.0.1:8080/todos/"+tPid, base.MakeToken(111), reply)
	if reply.Code != http.StatusForbidden || reply.Msg != BelongErr {
		t.Error("forbidden to delete other user's todo:", reply.Msg)
	}
	//normal case
	del("http://127.0.0.1:8080/todos/"+tPid, auth, reply)
	if reply.Code != 0 {
		t.Error("delete todo failed:", reply.Msg)
	}
}

func TestDeleteProject(t *testing.T) {
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	pid := cache.Get("delete_test_project_pid")
	reply := new(delProjectResp)
	del("http://127.0.0.1:8080/projects/"+pid, auth, reply)
	if reply.Code != 0 {
		t.Error("delete project failed:", reply.Msg)
	}
}
