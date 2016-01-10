package model

type Mission struct {
	Id        uint64
	StartTime uint64 //start timestamp of this action
	DeadLine  uint64 //end time
	Desc      string //description for the action
	Publisher User   //who published the mission
	Receivers []User //user list who received the mission
}
