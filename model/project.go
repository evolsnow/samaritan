package model

type Project struct {
	Id        int
	StartTime uint64 //start timestamp of this action
	DeadLine  uint64 //end time
	Desc      string //description for the action
	Color     [3]int //RGB mode
	Publisher User   //who published the mission
	Receivers []User //user list who received the mission
}
