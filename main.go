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
	confController := ConfigController{}
	authController := AuthInitialize(DB)
	errInit := authController.AuthMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatalf("Failed to initialize auth controller! %s", errInit)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/login", authController.AuthMiddleware.LoginHandler)
	r.POST("/logout", authController.AuthMiddleware.LogoutHandler)
	r.POST("/verify-stream", streamController.VerifyStream)
	r.POST("/stream-over", streamController.EndStream)
	streamGroup := r.Group("streams")
	streamGroup.Use(authController.AuthMiddleware.MiddlewareFunc())
	{
		streamGroup.GET("/", streamController.GetStreams)
		streamGroup.POST("/create-key", streamController.CreateKey)
		streamGroup.DELETE("/:id", streamController.DeleteStream)
	}
	nginxGroup := r.Group("/nginx")
	nginxGroup.Use(authController.AuthMiddleware.MiddlewareFunc())
	{
		nginxGroup.GET("/config", confController.GetConfiguration)
		nginxGroup.POST("/config", confController.UpdateConfiguration)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
