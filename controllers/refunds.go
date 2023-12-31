package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/BaseMax/real-time-notifications-nats-go/helpers"
	"github.com/BaseMax/real-time-notifications-nats-go/models"
	"github.com/BaseMax/real-time-notifications-nats-go/notifications"
	"github.com/BaseMax/real-time-notifications-nats-go/rabbitmq"
)

func AddRefund(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	order, dbErr := models.FindById[models.Order](uint(id))
	if dbErr != nil {
		return &dbErr.HttpErr
	}

	if order.Status != models.TASK_DONE {
		return echo.ErrNotAcceptable
	}

	refund := models.Refund{OrderID: uint(id)}
	if err := models.Create(&refund); err != nil {
		// GORM postgres driver doesn't have gorm.ErrForeignKeyViolated translation
		// I should hack
		if errors.Is(&err.HttpErr, echo.ErrInternalServerError) {
			return echo.ErrNotFound
		}
		return &err.HttpErr
	}

	user := helpers.GetLoggedinInfo(c)
	admin, dbErr := models.GetAdmin()
	if dbErr != nil {
		return &dbErr.HttpErr
	}

	activities := models.Activity{
		RecieverID: admin.ID,
		Title:      fmt.Sprintf("We have new refund from %s.", user.Username),
		Action:     models.ACTION_REFUND_NEW,
		Task:       models.AnyToTask(refund),
	}
	if err := notifications.Notify(activities); err != nil {
		return &err.HTTPError
	}

	if rabbitmq.RestartChannel() != nil {
		return echo.ErrInternalServerError
	}
	if err := rabbitmq.EnqueueTask(refund, rabbitmq.QUEUE_NAME_REFUNDS); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, refund)
}

func FetchRefund(c echo.Context) error {
	return GetModelByPreload[models.Refund](c, "id", "Order")
}

func FetchAllRefunds(c echo.Context) error {
	return GetAllModelsByPreload[models.Refund](c, "*", "Order.Products")
}

func EditRefund(c echo.Context) error {
	return EditModelById[models.Refund](c, "id")
}

func DeleteRefund(c echo.Context) error {
	return DeleteModelById[models.Refund](c, "id")
}

func FetchRefundStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	refund, dbErr := models.FindById[models.Refund](uint(id))
	if dbErr != nil {
		return &dbErr.HttpErr
	}
	return c.JSON(http.StatusOK, map[string]any{"status": refund.Status})
}

func CancelRefund(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	dbErr := models.UpdateById[models.Refund](uint(id), &models.Refund{Status: models.TASK_CANCELED})
	if dbErr != nil {
		return &dbErr.HttpErr
	}
	return c.JSON(http.StatusOK, map[string]any{"status": models.TASK_CANCELED})
}

func GetFirstQueuedRefund(c echo.Context) error {
	return ProcessFirstQueuedTask[models.Refund](c, rabbitmq.QUEUE_NAME_REFUNDS,
		models.TASK_BROWSE, models.TASK_BROWSE, "Order.Products")
}

func DoneFirstQueuedRefund(c echo.Context) error {
	return ProcessFirstQueuedTask[models.Refund](c, rabbitmq.QUEUE_NAME_REFUNDS,
		models.ACTION_REFUND_DONE, models.TASK_DONE, "Order.Products")
}

func CancelFirstQueuedRefund(c echo.Context) error {
	return ProcessFirstQueuedTask[models.Refund](c, rabbitmq.QUEUE_NAME_REFUNDS,
		models.ACTION_REFUND_CANCELED, models.TASK_CANCELED, "Order.Products")
}
