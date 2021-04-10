package main

import (
	"crypto/sha256"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	database *DatabaseUtility
}

func (a *AuthController) Login(ctx *gin.Context) {
	username := ctx.DefaultPostForm("username", "")
	password := hashPassword(ctx.DefaultPostForm("password", ""))

	rows, err := a.database.Login(username, password)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest, gin.H{
				SUCCESS_KEY:          false,
				RESPONSE_MESSAGE_KEY: "Login failed!",
			},
		)
	}
	defer rows.Close()
	accumulator := 0
	for rows.Next() {
		if accumulator > 0 {
			// We should only have one row to process.
			break
		}
		var id int
		err := rows.Scan(&id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{SUCCESS_KEY: false, RESPONSE_MESSAGE_KEY: err.Error()})
			return
		}
		jwt, err := createJwt(id)
		if err != nil {

		}
		accumulator++
	}
}

func hashPassword(password string) string {
	digest := sha256.New()
	digest.Write([]byte(password))
	return string(digest.Sum(nil))
}

func createJwt(userId int) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
