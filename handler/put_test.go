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

func put(reqURL string, src []byte, ds interface{}) {
	var t testing.T
	//reqURL = url.QueryEscape(reqURL)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", reqURL, bytes.NewReader(src))
	req.Header.Set("Content-Type", "application/json")
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
	src, _ := json.Marshal(req)
	reply := new(putPasswordResp)
	put("http://127.0.0.1:8080/users/password/gsc1215225@gmail.com", src, reply)
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
	src, _ = json.Marshal(req)
	cache.Set("gsc1215225@gmail.com:code", "123456", time.Minute*5)
	put("http://127.0.0.1:8080/users/password/gsc1215225@gmail.com", src, reply)
	if reply.Msg != CodeMismatchErr {
		t.Error("code mismatch")
	}

	req.VerifyCode = "123456"
	src, _ = json.Marshal(req)
	cache.Set("gsc@gmail.com:code", "123456", time.Minute*5)
	put("http://127.0.0.1:8080/users/password/gsc@gmail.com", src, reply)
	if reply.Msg != NotRegisteredErr {
		t.Error("not registed")
	}
}
