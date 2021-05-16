package main

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	db *DatabaseUtility
}

type CreateUserRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	IsAdmin  bool   `form:"isAdmin" json:"isAdmin" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
}

func (uc *UsersController) CreateUser(ctx *gin.Context) {
	jwtClaims := jwt.ExtractClaims(ctx)
	var newUser CreateUserRequest
	ctx.BindJSON(&newUser)
	if jwtClaims["isAdmin"].(bool) {
		insertedUser, err := uc.db.CreateUser(newUser.Username, newUser.Password, boolToInt(newUser.IsAdmin))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: "Failed to insert new user!",
			})
			return
		}
		lastInsertId, err2 := insertedUser.LastInsertId()
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: "Failed to get new user's id!",
			})
			return
		}
		newUser, err3 := uc.db.GetUser(lastInsertId)
		if err3 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: "Failed to get new user!",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY:          true,
			RESPONSE_MESSAGE_KEY: "Successfully inserted new user",
			"newUser":            newUser,
		})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "You don't have permission to create users!",
	})
	return
}

func (uc *UsersController) GetUsers(ctx *gin.Context) {
	jwtClaims := jwt.ExtractClaims(ctx)
	if jwtClaims["isAdmin"].(bool) {
		users, err := uc.db.GetUsers()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: "Failed to get users!",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY: true,
			"users":     users,
		})
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "You don't have permission to list users!",
	})
}
