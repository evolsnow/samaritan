package model

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"testing"
)

func TestGetOwner(t *testing.T) {
	tPid := cache.Get("delete_test_todo_pid")
	tid := dbms.ReadTodoId(tPid)
	td := Todo{Id: tid}
	u := td.GetOwner()
	if u == nil {
		t.Error("failed to get owner")
	}
}
