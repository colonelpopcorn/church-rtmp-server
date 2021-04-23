package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	database       *DatabaseUtility
	AuthMiddleware *jwt.GinJWTMiddleware
}

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type AuthorizedUser struct {
	UserId  float64
	IsAdmin bool
}

const userInfoKey = "userInfo"
const identityKey = "userId"

func AuthInitialize(db *DatabaseUtility) *AuthController {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {

			claims, _ := c.Get(userInfoKey)
			isAdmin := claims.(*AuthorizedUser).IsAdmin
			routes := []gin.H{
				{"path": "/config-editor", "name": "Configuration Editor"},
				{"path": "/user-manager", "name": "User Manager"},
			}
			nonAdminResponse := gin.H{
				"code":    http.StatusOK,
				"token":   token,
				"expire":  t.Format(time.RFC3339),
				"isAdmin": isAdmin,
			}
			adminResponse := gin.H{
				"code":    http.StatusOK,
				"token":   token,
				"expire":  t.Format(time.RFC3339),
				"isAdmin": isAdmin,
				"routes":  routes,
			}
			if isAdmin {
				c.JSON(http.StatusOK, adminResponse)
			} else {
				c.JSON(http.StatusOK, nonAdminResponse)
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AuthorizedUser); ok {
				return jwt.MapClaims{
					identityKey: v.UserId,
					"isAdmin":   v.IsAdmin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &AuthorizedUser{
				UserId:  claims[identityKey].(float64),
				IsAdmin: claims["isAdmin"].(bool),
			}
		},
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			var loginVals User
			if err := ctx.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			rows, err := db.Login(username, password)
			if err != nil {
				ctx.JSON(
					http.StatusBadRequest, gin.H{
						SUCCESS_KEY:          false,
						RESPONSE_MESSAGE_KEY: "Login failed!",
					},
				)
			}
			defer rows.Close()
			for rows.Next() {
				var id float64
				var isAdmin bool
				var hash string
				err := rows.Scan(&id, &isAdmin, &hash)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{SUCCESS_KEY: false, RESPONSE_MESSAGE_KEY: err.Error()})
					return nil, jwt.ErrFailedAuthentication
				}
				compare := CheckPasswordHash(password, hash)
				if compare {
					authedUser := &AuthorizedUser{
						UserId:  id,
						IsAdmin: isAdmin,
					}
					ctx.Set(userInfoKey, authedUser)
					log.Println(authedUser)
					return authedUser, nil
				}
			}
			return nil, errors.New("Illegal state, we shouldn't get here!")
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(*AuthorizedUser)
			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Printf("Error initializing jwt middleware! %s", err)
		return &AuthController{database: db}
	}
	return &AuthController{database: db, AuthMiddleware: authMiddleware}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
