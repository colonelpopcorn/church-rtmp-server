package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func TestUsersController(t *testing.T) {
	uc := UsersController{db: DbInitialize()}
	handler := ginHandler(uc)
	gfight := gofight.New()
	var authToken string
	gfight.POST("/login").
		SetJSON(gofight.D{"username": "admin", "password": "admin"}).
		Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var message gin.H
			json.Unmarshal(r.Body.Bytes(), &message)
			assert.NotEqual(t, message["token"], nil)
			authToken = message["token"].(string)
		})
	t.Run("Test getting users", func(t *testing.T) {
		gfight.
			GET("/users/").
			SetHeader(gofight.H{"Authorization": fmt.Sprintf("Bearer %s", authToken)}).
			Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				var message gin.H
				json.Unmarshal(r.Body.Bytes(), &message)
				t.Log(message)
				assert.Equal(t, message[SUCCESS_KEY], true)
			})
	})
	t.Run("Test creating and deleting user", func(t *testing.T) {
		gfight.
			POST("/users/create").
			SetJSON(gofight.D{
				"username": "notAdmin" + generateGUID(),
				"password": "billyBob92",
				"isAdmin":  false,
			}).
			Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				var message gin.H
				json.Unmarshal(r.Body.Bytes(), &message)
				t.Log(message)
				assert.Equal(t, message[SUCCESS_KEY], true)
			})
		createdUserId, _ := getLastInsertedUser(uc)
		gfight.
			DELETE(fmt.Sprintf("/users/%d", createdUserId)).
			Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				var message gin.H
				json.Unmarshal(r.Body.Bytes(), &message)
				t.Log(message)
				assert.Equal(t, message[SUCCESS_KEY], true)
			})
	})
	t.Run("Test updating user password", func(t *testing.T) {
		gfight.
			POST("/users/create").
			SetJSON(gofight.D{
				"username": "notAdmin" + generateGUID(),
				"password": "billyBob92",
				"isAdmin":  false,
			}).
			Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				var message gin.H
				json.Unmarshal(r.Body.Bytes(), &message)
				t.Log(message)
				assert.Equal(t, message[SUCCESS_KEY], true)
			})
		createdUserId, _ := getLastInsertedUser(uc)
		gfight.
			POST("/users/update").
			SetJSON(gofight.D{
				"userId":      createdUserId,
				"oldPassword": "billyBob92",
				"newPassword": "billyBob93" + generateGUID(),
			}).
			Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				var message gin.H
				json.Unmarshal(r.Body.Bytes(), &message)
				t.Log(message)
				assert.Equal(t, message[SUCCESS_KEY], true)
			})
		gfight.
			DELETE(fmt.Sprintf("/users/%d", createdUserId)).
			Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				var message gin.H
				json.Unmarshal(r.Body.Bytes(), &message)
				t.Log(message)
				assert.Equal(t, message[SUCCESS_KEY], true)
			})
	})
}

func ginHandler(uc UsersController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	auth := getAuthMiddleWare()

	r.POST("/login", auth.LoginHandler)
	r.POST("/logout", auth.LogoutHandler)
	// test token in path
	r.GET("/g/:token/refresh_token", auth.RefreshHandler)

	authGroup := r.Group("/auth")
	// Refresh time can be longer than token timeout
	authGroup.GET("/refresh_token", auth.RefreshHandler)
	group := r.Group("/users")
	group.Use(auth.MiddlewareFunc())
	{
		group.GET("/", uc.GetUsers)
		group.POST("/create", uc.CreateUser)
		group.DELETE("/:id", uc.DeleteUser)
		group.POST("/update", uc.UpdateUserPassword)
	}

	return r
}

func getAuthMiddleWare() jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "test zone",
		Key:   []byte("secret key"),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// Set custom claim, to be checked in Authorizator method
			return jwt.MapClaims{"testkey": "testval", "exp": 0, "isAdmin": true}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if binderr := c.ShouldBind(&loginVals); binderr != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			if userID == "admin" && password == "admin" {
				return userID, nil
			}
			return "", jwt.ErrFailedAuthentication
		},
		Authorizator: func(user interface{}, c *gin.Context) bool {
			return true
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			cookie, err := c.Cookie("jwt")
			if err != nil {
				log.Println(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"token":   token,
				"expire":  t.Format(time.RFC3339),
				"message": "login successfully",
				"cookie":  cookie,
			})
		},
		TimeFunc: func() time.Time { return time.Now().Add(time.Duration(5) * time.Minute) },
	})
	if err != nil {
		panic(err)
	}
	return *authMiddleware
}

func getLastInsertedUser(uc UsersController) (int64, error) {
	var createdUserId int64
	rows, err := uc.db.dbContext.Query("SELECT MAX(id) from users;")
	if err != nil {
		return 0, errors.New("could not get max id")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&createdUserId)
		if err != nil {
			return 0, errors.New("could not get max id")
		}
	}
	return createdUserId, nil
}
