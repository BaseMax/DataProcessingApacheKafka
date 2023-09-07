package models

const (
	TASK_INPROGRESS = "IN_PROGRESS"
	TASK_DONE       = "DONE"
	TASK_CANCELED   = "CANCELED"
	TASK_BROWSE     = "BROWSE"

	ACTION_REGISTER        = "REGISTER"
	ACTION_ORDER_NEW       = "ORDER_NEW"
	ACTION_ORDER_DONE      = "ORDER_DONE"
	ACTION_ORDER_CANCELED  = "ORDER_CANCELED"
	ACTION_REFUND_NEW      = "REFUND_NEW"
	ACTION_REFUND_DONE     = "REFUND_DONE"
	ACTION_REFUND_CANCELED = "REFUND_CANCELED"
)

var ACTIONS = []string{
	ACTION_REGISTER, ACTION_ORDER_NEW,
	ACTION_ORDER_DONE, ACTION_ORDER_CANCELED,
	ACTION_REFUND_NEW, ACTION_REFUND_DONE,
	ACTION_REFUND_CANCELED,
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
