package model

type Todo struct {
	Id        int
	StartTime uint64  //start timestamp of this action
	DeadLine  uint64  //end time
	Desc      string  //description for the action
	Owner     User    //whose
	Status    int     //0:not begin, 1:ongoing, 2: overdue, 3:accomplished
	Belong    Project //belong to which project
}

//TodoList
//	类型	备注
//tableID	int	主键
//dayID	int	查询关键词，如20160223
//startTime	long long	开始时间戳
//endTime	long long	结束时间戳
//thing	TodoThing	事件，TodoThing表见下图
//doneType	int	0 未开始 1进行中 2过期 3已完成
//
//TodoThing	类型	备注
//
//thingStr	string	事件描述字符串
//Images	array	图片数组
//thingType	TodoThingType	事件类型，表见下图
//
//TodoThingType	类型	备注
//
//typeId	int
//typeStr	string	类型描述字符串
//typeRed	int	RGB值
//typeGreen	int
//typeBlue	int
