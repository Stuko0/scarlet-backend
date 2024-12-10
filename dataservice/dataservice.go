package dataservice

import (
	"scarlet_backend/dataservice/firedata"
	"scarlet_backend/dataservice/newdata"
	"scarlet_backend/dataservice/taskdata"
	"scarlet_backend/dataservice/teamdata"
	"scarlet_backend/dataservice/userdata"
)

type DataService struct {
	FireData   firedata.FireDataMongo
	TaskData   taskdata.TaskDataPostgres
	UserData   userdata.UserDataPostgres
	MemberData teamdata.MemberDataPostgres
	TeamData   teamdata.TeamDataPostgres
	NewsData   newdata.NewsDataPostgres
}
