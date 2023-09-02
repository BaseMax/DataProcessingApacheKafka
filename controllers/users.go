package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/BaseMax/real-time-notifications-nats-go/models"
)

var EXPTIME = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30))

func createJwtToken(id uint, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprint(id),
		Issuer:    username,
		ExpiresAt: EXPTIME,
	})
	bearer, _ := token.SignedString(GetJwtSecret())
	return bearer
}

func getLoggedinInfo(c echo.Context) (uint, string) {
	bearer := c.Request().Header.Get("Authorization")
	token, _, _ := new(jwt.Parser).ParseUnverified(bearer[len("Bearer "):], jwt.MapClaims{})
	claims := token.Claims.(jwt.MapClaims)

	username := claims["iss"].(string)
	id, _ := strconv.Atoi(claims["jti"].(string))
	return uint(id), username
}

func Register(c echo.Context) error {
	var user models.User
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return echo.ErrBadRequest
	}
	user.Password = models.HashPassword(user.Password)
	if err := models.Create(&user); err != nil {
		return &err.HttpErr
	}
	bearer := createJwtToken(user.ID, user.Username)
	return c.JSON(http.StatusOK, map[string]any{"bearer": bearer})
}

func Login(c echo.Context) error {
	var user models.User
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return echo.ErrBadRequest
	}
	user.Password = models.HashPassword(user.Password)
	if err := models.Login(&user); err != nil {
		return &err.HttpErr
	}
	bearer := createJwtToken(user.ID, user.Username)
	return c.JSON(http.StatusOK, map[string]any{"bearer": bearer})
}

func Refresh(c echo.Context) error {
	bearer := createJwtToken(getLoggedinInfo(c))
	return c.JSON(http.StatusOK, map[string]any{"bearer": bearer})
}

func FetchUser(c echo.Context) error {
	return nil
}

func FetchAllUsers(c echo.Context) error {
	return nil
}

func DeleteUser(c echo.Context) error {
	return nil
}

func EditUser(c echo.Context) error {
	return nil
}
