package main

import (
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	database *DatabaseUtility
}

func (a *AuthController) Login(_ *gin.Context) {
}
