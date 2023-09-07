package models

const (
	TASK_INPROGRESS = "IN_PROGRESS"
	TASK_DONE       = "DONE"
	TASK_CANCELED   = "CANCELED"
	TASK_BROWSE     = "BROWSE"

	ACTION_REGISTER   = "REGISTER"
	ACTION_NEW_RECORD = "NEW_RECORD"
)

var ACTIONS = []string{
	TASK_INPROGRESS, TASK_DONE,
	TASK_CANCELED, TASK_BROWSE,
	ACTION_REGISTER, ACTION_NEW_RECORD,
}

type Activity struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RecieverID uint   `gorm:"not null" json:"reciever_id"`
	Title      string `gorm:"not null" json:"title"`
	Action     string `gorm:"not null" json:"action"`

	Task Task `gorm:"-" json:"task,omitempty"`
}

func GetActivitiesByUserId(id uint) (actitivies *[]Activity, herr *DbErr) {
	err := db.Where(&Activity{RecieverID: id}).Find(&actitivies)
	return actitivies, errGormToHttp(err)
}

func SeenActivities(id uint) *DbErr {
	err := db.Where(&Activity{RecieverID: id}).Delete(&Activity{})
	return errGormToHttp(err)
}
