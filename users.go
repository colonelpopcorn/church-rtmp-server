package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	UserId      int64  `form:"userId" json:"userId" binding:"required"`
}

func (uc *UsersController) CreateUser(ctx *gin.Context) {
	jwtClaims := jwt.ExtractClaims(ctx)
	if jwtClaims["isAdmin"].(bool) {
		var newUser CreateUserRequest
		ctx.BindJSON(&newUser)
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
		returnUser := User{UserId: lastInsertId, IsAdmin: newUser.IsAdmin, Username: newUser.Username}
		ctx.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY:          true,
			RESPONSE_MESSAGE_KEY: "Successfully inserted new user",
			"newUser":            returnUser,
		})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "You don't have permission to create users!",
	})
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
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "You don't have permission to list users!",
	})
}

func (uc *UsersController) DeleteUser(ctx *gin.Context) {
	jwtClaims := jwt.ExtractClaims(ctx)
	if jwtClaims["isAdmin"].(bool) {
		userIdStr := ctx.Param("id")
		userId, _ := strconv.ParseInt(userIdStr, 10, 64)
		_, err := uc.db.DeleteUser(userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: fmt.Sprintf("Failed to delete existing record, %s", err.Error()),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY:          true,
			RESPONSE_MESSAGE_KEY: "Deleted user!",
		})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "You don't have permission to delete users!",
	})
}

func (uc *UsersController) UpdateUserPassword(ctx *gin.Context) {
	jwtClaims := jwt.ExtractClaims(ctx)
	if jwtClaims["isAdmin"].(bool) {
		var userUpdateReq ChangePasswordRequest
		ctx.BindJSON(&userUpdateReq)
		_, err := uc.db.UpdateUserPassword(userUpdateReq.UserId, userUpdateReq.OldPassword, userUpdateReq.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: fmt.Sprintf("Failed to delete existing record, %s", err.Error()),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			SUCCESS_KEY:          true,
			RESPONSE_MESSAGE_KEY: "Successfully changed user password!",
		})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		SUCCESS_KEY:          false,
		RESPONSE_MESSAGE_KEY: "You don't have permission to change user passwords!",
	})
}
