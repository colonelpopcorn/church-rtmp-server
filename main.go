package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type SessionToken struct {
	Token string `json:"token"`
}

func main() {
	DB := DbInitialize()
	defer DB.CloseDb()
	streamController := StreamController{DB}
	confController := ConfigController{DB}
	authController := AuthInitialize(DB)
	errInit := authController.AuthMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatalf("Failed to initialize auth controller! %s", errInit)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/login", authController.AuthMiddleware.LoginHandler)
	r.POST("/verify-stream", streamController.VerifyStream)
	r.GET("/streams", streamController.GetStreams)
	r.POST("/create-key", streamController.CreateKey)
	r.DELETE("/streams/:id", streamController.DeleteStream)
	r.POST("/stream-over", streamController.EndStream)
	r.GET("/nginx-conf", confController.GetConfiguration)
	r.POST("/nginx-conf", confController.UpdateConfiguration)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
