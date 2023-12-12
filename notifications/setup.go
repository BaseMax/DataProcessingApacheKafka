package notifications

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/segmentio/kafka-go"
)

var (
	nConn *nats.Conn
	kConn *kafka.Conn
)

func InitNats() error {
	var err error
	nConn, err = nats.Connect(GetNatsURL())
	return err
}

func GetNatsConn() *nats.Conn {
	return nConn
}

func InitKafka() error {
	var err error

	conf, err := GetKafkaConfig()
	if err != nil {
		return err
	}

	kConn, err = kafka.DialLeader(context.Background(),
		conf.Protocol, conf.Server, conf.Topic, conf.Partition)
	return err
}

func GetKafkaConn() *kafka.Conn {
	return kConn
}
