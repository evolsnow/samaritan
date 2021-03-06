package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/common/base"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/model"
	"net/http"
	"testing"
)

func get(reqURL, auth string, ds interface{}) {
	var t testing.T
	//reqURL = url.QueryEscape(reqURL)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("http get err")
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(ds)
	if err != nil {
		t.Error(err)
	}
}

func TestSamIdStatus(t *testing.T) {

	reply := new(samIdStatusResp)

	dbms.DeleteSamId("testevol")
	get("http://127.0.0.1:8080/samIds/testevol", "", reply)
	if reply.Code != 0 {
		t.Error("should be available")
	}

	dbms.UpdateSamIdSet("testevol")
	get("http://127.0.0.1:8080/samIds/testevol", "", reply)
	if reply.Code == 0 || reply.Msg != ExistErr {
		t.Error("should be unavailable")
	}

	get(`http://127.0.0.1:8080/samIds/*!1234`, "", reply)
	if reply.Msg != CharsetErr {
		t.Error("illegal charset")
	}

	get("http://127.0.0.1:8080/samIds/abc", "", reply)
	if reply.Msg != LengthErr {
		t.Error("illegal length")
	}

}

func TestUserProjectList(t *testing.T) {

	reply := new(userProjectsResp)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	p := &model.Project{
		CreatorId: uid,
		Name:      "pj name",
		Desc:      "pj desc",
	}
	p.Save()
	auth := base.MakeToken(uid)
	get("http://127.0.0.1:8080/projects", "", reply)
	if reply.Code == 0 {
		t.Error("should be unauthorized")
	}
	get("http://127.0.0.1:8080/projects", auth, reply)
	if reply.Code != 0 || len(reply.Np) < 1 {
		t.Error("failed to get projects")
	}
	get("http://127.0.0.1:8080/projects?type=joined", auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get joined projects")
	}
	get("http://127.0.0.1:8080/projects?type=created", auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get created projects")
	}
	get("http://127.0.0.1:8080/projects?type=fake", auth, reply)
	if reply.Msg != UnknownTypeErr {
		t.Error("should be unknow type")
	}
}

func TestUserSearch(t *testing.T) {
	reply := new(userSearchResp)
	get("http://127.0.0.1:8080/users/pubInfo/mail@fake.com", "", reply)
	if reply.Msg != UserNotExistErr {
		t.Error("user not exist")
	}
	reply = new(userSearchResp)
	get("http://127.0.0.1:8080/users/pubInfo/gsc1215225@gmail.com", "", reply)
	if reply.Code != 0 || reply.Id == "" {
		t.Error("failed to serach user")
	}
}

func TestProjectMissionList(t *testing.T) {

	reply := new(projectMissionsResp)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	p := &model.Project{
		CreatorId: uid,
		Name:      "pj2 name",
		Desc:      "pj2 desc",
	}
	p.Save()
	m := &model.Mission{
		PublisherId:   uid,
		ProjectId:     p.Id,
		Name:          "ms name",
		Desc:          "ms desc",
		Deadline:      147258369,
		CompletionNum: 70,
	}
	m.Save()
	auth := base.MakeToken(uid)
	get("http://127.0.0.1:8080/projectMissions/"+p.Pid, "", reply)
	if reply.Code == 0 {
		t.Error("should be unauthorized")
	}
	get("http://127.0.0.1:8080/projectMissions/"+p.Pid, auth, reply)
	if reply.Code != 0 || len(reply.Nm) < 1 {
		t.Error("failed to get missions:", reply.Msg)
	}
	get("http://127.0.0.1:8080/projectMissions/"+p.Pid, base.MakeToken(dbms.ReadUserIdWithIndex("snow@gmail.com", "mail")), reply)
	if reply.Code == 0 || reply.Msg != NotProjectMemberErr {
		t.Error("not member:", reply.Msg)
	}
}

func TestMissionCommentList(t *testing.T) {

	reply := new(missionCommentResp)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")

	m := &model.Mission{
		PublisherId:   uid,
		Name:          "ms name 3",
		Desc:          "ms desc 3",
		CompletionNum: 10,
	}
	m.Save()
	auth := base.MakeToken(uid)
	get("http://127.0.0.1:8080/comments/"+m.Pid, auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get comment")
	}
}

func TestMissionDetail(t *testing.T) {
	mPid := cache.Get("put_test_mission_pid")
	reply := new(missionDetailResp)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	get("http://127.0.0.1:8080/missions/"+mPid, auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get mission")
		t.FailNow()
	}
	if reply.PublisherId == "" {
		t.Error("failed to get mission publisher")
	}
	if len(reply.ReceiversName) == 0 {
		t.Error("faield to get mission receivers")
	}
}

func TestOfflineMsgs(t *testing.T) {
	reply := new(getOfflineMsgResp)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	get("http://127.0.0.1:8080/offlineMessages", auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get messages")
	}
	if len(reply.Msgs) > 0 && reply.Msgs[0].Id == "" {
		t.Error("got wrong chat id")
	}
}

func TestUserInfo(t *testing.T) {
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	reply := new(personalInfoResp)
	get("http://127.0.0.1:8080/users/personalInfo", auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get user personal info")
	}
	if reply.Id == "" {
		t.Error("failed to get user id info")
	}
}

func TestProjectDetail(t *testing.T) {
	pPid := cache.Get("put_test_project_pid")
	reply := new(projectDetailResp)
	uid := dbms.ReadUserIdWithIndex("gsc1215225@gmail.com", "mail")
	auth := base.MakeToken(uid)
	get("http://127.0.0.1:8080/projects/"+pPid, auth, reply)
	if reply.Code != 0 {
		t.Error("failed to get mission")
		t.FailNow()
	}
	if reply.CreatorId == "" {
		t.Error("failed to get project creator")
	}
	if len(reply.Members) == 0 {
		t.Error("faield to get project members")
	}
}
