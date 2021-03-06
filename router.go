package main

import (
	"github.com/evolsnow/httprouter"
	"github.com/evolsnow/samaritan/common/base"
	hd "github.com/evolsnow/samaritan/handler"
	mw "github.com/evolsnow/samaritan/middleware"
	"net/http"
)

//const BaseURL = "/api/1.0"

var r = httprouter.New()

type notAllowed struct{}
type notFound struct{}

func newRouter() *httprouter.Router {
	//base url for all requests
	//r.BaseURL = BaseURL

	//user-defined http error handler
	r.MethodNotAllowed = notAllowed{}
	r.NotFound = notFound{}

	//http get method
	r.GET("/samIds/:samId", hd.SamIdStatus)
	r.GET("/users/pubInfo/:user", hd.SearchUser)
	r.GET("/users/personalInfo", mw.Auth(hd.UserInfo))
	r.GET("/projects", mw.Auth(hd.UserProjectList))
	r.GET("/projects/:project", mw.Auth(hd.ProjectDetail))
	r.GET("/projectMissions/:project", mw.Auth(hd.ProjectMissionList))
	r.GET("/comments/:mission", mw.Auth(hd.MissionCommentList))
	r.GET("/missions/:mission", mw.Auth(hd.MissionDetail))
	r.GET("/offlineMessages", mw.Auth(hd.OfflineMsgs))
	r.GET("/qiniu/uploadTokens", mw.Auth(hd.MakeUploadToken))

	//http post method
	r.POST("/deviceTokens", mw.Auth(hd.NewDeviceToken))
	r.POST("/users", hd.NewUser)
	r.POST("/verifyCode/:source", hd.NewVerifyCode)
	r.POST("/accessToken", hd.NewAccessToken)
	r.POST("/todos", mw.Auth(hd.NewTodo))
	r.POST("/missions", mw.Auth(hd.NewMission))
	r.POST("/comments", mw.Auth(hd.NewComment))
	r.POST("/projects", mw.Auth(hd.NewProject))
	r.POST("/privateChats", mw.Auth(hd.NewPrivateChat))
	r.POST("/invitations/project", mw.Auth(hd.NewProjectInvitation))
	r.POST("/invitations/mission", mw.Auth(hd.NewMissionInvitation))

	//http put method
	r.PUT("/users/password/:identity", hd.UpdatePassword)
	r.PUT("/users/personalInfo", mw.Auth(hd.UpdateUserInfo))
	r.PUT("/todos/:todo", mw.Auth(hd.UpdateTodo))
	//r.PUT("/todos/pictures/:todo", mw.Auth(hd.UpdateTodoPics))
	r.PUT("/missions/pictures/:mission", mw.Auth(hd.UpdateMissionPics))
	r.PUT("/missions/status/:mission", mw.Auth((hd.UpdateMissionStatus)))
	r.PUT("/missions/accepted/:mission", mw.Auth(hd.AcceptMission))
	r.PUT("/projects/joined/:project", mw.Auth(hd.JoinProject))
	r.PUT("/chats/status/:chat", mw.Auth((hd.UpdateChatStatus)))

	//http delete method
	r.DELETE("/todos/:todo", mw.Auth(hd.DeleteTodo))
	r.DELETE("/missions/:mission", mw.Auth(hd.DeleteMission))
	r.DELETE("/projects/:project", mw.Auth(hd.DeleteProject))

	return r
}

// ServeHTTP makes the NAllowed implement the http.Handler interface.
func (notAllowed) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	base.MethodNAErr(w, http.StatusText(http.StatusMethodNotAllowed))
}

// ServeHTTP makes the NFound implement the http.Handler interface.
func (notFound) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	base.NotFoundErr(w, http.StatusText(http.StatusNotFound))
}
