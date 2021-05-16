package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateUser(t *testing.T) {
	r := gin.Default()

	uc := UsersController{}
	r.GET("/users", uc.GetUsers)
}

func getDbMock() (db IDatabaseUtility) {
}
