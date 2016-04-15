package model

import (
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"strconv"
	"time"
)

func init() {
	go func() {
		for {
			now := time.Now()
			// next 4am
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 4, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			//SyncMysql()
		}
	}()
}

const (
	UserInsert = "INSERT INTO user(redis_id, pid, sam_id, create_time, alias, name, phone, password, email, avatar, school, depart, grade, class, student_num) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	UserUpdate = "UPDATE user SET sam_id=?, alias=?, name=?, phone=?, password=?, email=?, avatar=?, school=?, depart=?, grade=?, class=?, student_num=? WHERE redis_id=?"
	UserDelete = "DELETE FROM user WHERE redis_id = ?"

	TodoInsert = "INSERT INTO todo(redis_id, pid, create_time, start_time, place, is_repeat, repeat_mode, all_day, describtion, remark, owner_id, done, finish_time, mission) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	TodoUpdate = "UPDATE todo SET start_time=?, place=?, is_repeat=?, repeat_mode=?, all_day=?, describtion=?, remark=?, done=?, finish_time=? WHERE redis_id=?"
	TodoDelete = "DELETE FROM todo WHERE redis_id = ?"

	MissionInsert = "INSERT INTO mission(redis_id, pid, create_time, name, describtion, publisher, completion_num, completed_time, deadline, project) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	MissionUpdate = "UPDATE mission SET name=?, describtion=?, publisher=?, completion_num=?, completed_time=?, deadine=? WHERE redis_id=?"
	MissionDelete = "DELETE FROM mission WHERE redis_id = ?"

	ProjectInsert = "INSERT INTO project(redis_id, pid, create_time, name, describtion, background_pic, creator, private) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	ProjectUpdate = "UPDATE project SET name=?, describtion=?, background_pic=?, private=? WHERE redis_id=?"
	ProjectDelete = "DELETE FROM project WHERE redis_id = ?"
)

func CreateUserMysql(u User) {
	stmt, err := dbms.DB.Prepare(UserInsert)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(u.Id, u.Pid, u.SamId, u.CreateTime, u.Alias, u.Name, u.Phone, u.Password, u.Email, u.Avatar, u.School, u.Department, u.Grade, u.Class, u.StudentNum)
}

func UpdateUserMysql(u User) {
	stmt, err := dbms.DB.Prepare(UserUpdate)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(u.SamId, u.Alias, u.Name, u.Phone, u.Password, u.Email, u.Avatar, u.School, u.Department, u.Grade, u.Class, u.StudentNum, u.Id)

}

func DeleteUserMysql(uid int) {
	stmt, err := dbms.DB.Prepare(UserDelete)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(uid)
}

func CreateTodoMysql(t Todo) {
	stmt, err := dbms.DB.Prepare(TodoInsert)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(t.Id, t.Pid, t.CreateTime, t.StartTime, t.Place, t.Repeat, t.RepeatMode, t.AllDay, t.Desc, t.Remark, t.OwnerId, t.Done, t.FinishTime, t.MissionId)
}

func UpdateTodoMysql(t Todo) {
	stmt, err := dbms.DB.Prepare(TodoUpdate)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(t.StartTime, t.Place, t.Repeat, t.RepeatMode, t.AllDay, t.Desc, t.Remark, t.Done, t.FinishTime, t.Id)
}

func DeleteTodoMysql(tid int) {
	stmt, err := dbms.DB.Prepare(TodoDelete)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(tid)
}

func CreateMissionMysql(m Mission) {
	stmt, err := dbms.DB.Prepare(MissionInsert)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(m.Id, m.Pid, m.CreateTime, m.Name, m.Desc, m.PublisherId, m.CompletionNum, m.CompletedTime, m.Deadline, m.ProjectId)
}

func UpdateMissionMysql(m Mission) {
	stmt, err := dbms.DB.Prepare(MissionUpdate)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(m.Name, m.Desc, m.PublisherId, m.CompletionNum, m.CompletedTime, m.Deadline, m.Id)
}

func DeleteMissionMysql(mid int) {
	stmt, err := dbms.DB.Prepare(MissionDelete)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(mid)
}

func CreateProjectMysql(p Project) {
	stmt, err := dbms.DB.Prepare(ProjectInsert)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(p.Id, p.Pid, p.CreateTime, p.Name, p.Desc, p.BackgroundPic, p.CreatorId, p.Private)
}

func UpdateProjectMysql(p Project) {
	stmt, err := dbms.DB.Prepare(ProjectUpdate)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(p.Name, p.Desc, p.BackgroundPic, p.Private, p.Id)
}

func DeleteProjectMysql(pid int) {
	stmt, err := dbms.DB.Prepare(ProjectDelete)
	if err != nil {
		log.Error(err)
	}
	stmt.Exec(pid)
}

func SyncMysql() {
	syncUser()
	syncTodo()
	syncMission()
	syncProject()
}

func syncUser() {
	tmp, _ := dbms.Get("autoIncrUser")
	total, _ := strconv.Atoi(tmp)
	delete, _ := dbms.DB.Prepare(UserDelete)
	update, _ := dbms.DB.Prepare(UserUpdate)
	for i := 1; i <= total; i++ {
		u, _ := readUserWithId(i)
		if u.Id == 0 {
			//delete
			delete.Exec(i)
		} else {
			//update
			update.Exec(u.SamId, u.Alias, u.Name, u.Phone, u.Password, u.Email, u.Avatar, u.School, u.Department, u.Grade, u.Class, u.StudentNum, u.Id)
		}
	}
	log.Info("user synced")
}

func syncTodo() {
	tmp, _ := dbms.Get("autoIncrTodo")
	total, _ := strconv.Atoi(tmp)
	delete, _ := dbms.DB.Prepare(TodoDelete)
	update, _ := dbms.DB.Prepare(TodoUpdate)
	for i := 1; i <= total; i++ {
		t, _ := readTodoWithId(i)
		if t.Id == 0 {
			//delete
			delete.Exec(i)
		} else {
			//update
			update.Exec(t.StartTime, t.Place, t.Repeat, t.RepeatMode, t.AllDay, t.Desc, t.Remark, t.Done, t.FinishTime, t.Id)
		}
	}
	log.Info("todo synced")
}

func syncMission() {
	tmp, _ := dbms.Get("autoIncrMission")
	total, _ := strconv.Atoi(tmp)
	delete, _ := dbms.DB.Prepare(MissionDelete)
	update, _ := dbms.DB.Prepare(MissionUpdate)
	for i := 1; i <= total; i++ {
		m, _ := readMissionWithId(i)
		if m.Id == 0 {
			//delete
			delete.Exec(i)
		} else {
			//update
			update.Exec(m.Name, m.Desc, m.PublisherId, m.CompletionNum, m.CompletedTime, m.Deadline, m.Id)
		}
	}
	log.Info("mission synced")
}

func syncProject() {
	tmp, _ := dbms.Get("autoIncrProject")
	total, _ := strconv.Atoi(tmp)
	delete, _ := dbms.DB.Prepare(ProjectDelete)
	update, _ := dbms.DB.Prepare(ProjectUpdate)
	for i := 1; i <= total; i++ {
		p, _ := readProjectWithId(i)
		if p.Id == 0 {
			//delete
			delete.Exec(i)
		} else {
			//update
			update.Exec(p.Name, p.Desc, p.BackgroundPic, p.Private, p.Id)
		}
	}
	log.Info("project synced")
}

func Test() {
	var name string
	row := dbms.DB.QueryRow("select name from user where id = ?", 1)
	err := row.Scan(&name)
	if err != nil {
		log.Error(err)
	}
	log.Println(name)
}
