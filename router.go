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

type NotAllowed struct{}
type NotFound struct{}

func newRouter() *httprouter.Router {
	//base url for all requests
	//r.BaseURL = BaseURL

	//user-defined http error handler
	r.MethodNotAllowed = NotAllowed{}
	r.NotFound = NotFound{}

	//http get method
	r.GET("/samIds/:samId", hd.SamIdStatus)
	r.GET("/projects", mw.Auth(hd.UserProjectList))
	r.GET("/projects/missions/:project", mw.Auth(hd.ProjectMissionList))
	r.GET("/missions/comments/:mission", mw.Auth(hd.MissionCommentList))
	r.GET("/users/:user", hd.SearchUser)

	//http post method
	r.POST("/users", hd.NewUser)
	r.POST("/verifyCode/:source", hd.NewVerifyCode)
	r.POST("/accessToken", hd.NewAccessToken)
	r.POST("/todos", mw.Auth(hd.NewTodo))
	r.POST("/projects", mw.Auth(hd.NewProject))
	r.POST("/privateChats", mw.Auth(hd.NewPrivateChat))
	r.POST("/invitations", mw.Auth(hd.NewInvitation))
	//todo upload device token

	//http put method
	r.PUT("/users/password/:identity", hd.UpdatePassword)
	r.PUT("/todos/:todo", mw.Auth(hd.UpdateTodo))
	r.PUT("/missions/status/:mission", mw.Auth((hd.UpdateMissionStatus)))

	//http delete method
	r.DELETE("/todos/:todo", mw.Auth(hd.DeleteTodo))
	r.DELETE("/projects/:project", mw.Auth(hd.DeleteProject))

	//test
	r.GET("/test", hd.Test)
	r.POST("/hi", hd.Hi)
	r.POST("/hia/:userId", mw.Auth(hd.Hi))
	r.GET("/set", hd.SetJwt)
	r.GET("/pm", hd.Pm)
	r.GET("/pm2", hd.Pm2)
	r.GET("/ab", hd.Ab)

	return r
}

// ServeHTTP makes the NAllowed implement the http.Handler interface.
func (NotAllowed) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	base.MethodNAErr(w, http.StatusText(http.StatusMethodNotAllowed))
}

// ServeHTTP makes the NFound implement the http.Handler interface.
func (NotFound) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	base.NotFoundErr(w, http.StatusText(http.StatusNotFound))
}
