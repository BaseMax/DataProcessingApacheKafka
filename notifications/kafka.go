package notifications

import (
	"encoding/json"

	"github.com/BaseMax/real-time-notifications-nats-go/models"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

func NotifyKafka(activity models.Activity) *NotifErr {
	data, _ := json.Marshal(activity)
	_, err := kConn.WriteMessages(kafka.Message{Value: data})
	if err != nil {
		return &NotifErr{Err: err, HTTPError: *echo.ErrInternalServerError}
	}
	return nil
}
