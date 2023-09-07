package notifications

import (
	"fmt"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

type NotifErr struct {
	Err       error
	HTTPError echo.HTTPError
}

func (e *NotifErr) Error() string { return e.Err.Error() }

func (e *NotifErr) Unwrap() error { return e.Err }

func GetNatsURL() string {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return nats.DefaultURL
	}
	return url
}

func CreateSubject(id uint) string {
	return fmt.Sprintf("notify.%d", id)
}

type KafkaConfig struct {
	Protocol  string
	Server    string
	Topic     string
	Partition int
}

func GetKafkaConfig() (*KafkaConfig, error) {
	partition, err := strconv.Atoi(os.Getenv("KAFKA_PARTITION"))
	if err != nil {
		return nil, err
	}
	return &KafkaConfig{
		Protocol:  "tcp",
		Server:    os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
		Topic:     os.Getenv("KAFKA_NOTIFICATION_TOPIC"),
		Partition: partition,
	}, nil
}
