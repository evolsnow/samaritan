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
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
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
	if reply.Code == 200 {
		t.Error("unauthorized to delete this todo")
	}

	del("http://127.0.0.1:8080/todos/"+tPid, auth, reply)
	if reply.Code != 200 {
		t.Error("delete to do failed:", reply.Msg)
	}
}
