package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/BaseMax/real-time-notifications-nats-go/models"
)

func ProcessFirstTask[T any](queueName, taskAction, preload string) (model *T, err error) {
	tasks, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	for {
		var modelOnQ T
		var interfaceModel any

		select {
		case task := <-tasks:
			if err := json.Unmarshal(task.Body, &modelOnQ); err != nil {
				log.Println(err)
			}

			interfaceModel = modelOnQ
			model, ok := interfaceModel.(models.Task)
			if !ok {
				log.Println("type is not Task interface")
			}

			modelOnDB, dbErr := models.FindByIdPreload[T](model.GetID(), preload)
			if dbErr != nil {
				log.Println(dbErr.Err)
			}

			interfaceModel = modelOnDB
			model, ok = interfaceModel.(models.Task)
			if !ok {
				log.Println("type is not Task interface")
			}

			if model.GetStatus() != models.TASK_INPROGRESS {
				task.Ack(false)
			} else {
				if taskAction == models.TASK_DONE || taskAction == models.TASK_CANCELED {
					models.UpdateStatus[T](model.GetID(), taskAction)
					task.Ack(false)
				} else {
					task.Nack(false, true)
				}
				return modelOnDB, nil
			}

		case <-time.After(time.Second / 2):
			return nil, nil
		}
	}
}
