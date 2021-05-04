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
	app := gin.Default()
	app.Use(cors.Default())
	app.POST("/login", authController.AuthMiddleware.LoginHandler)
	app.POST("/logout", authController.AuthMiddleware.LogoutHandler)
	app.POST("/verify-stream", streamController.VerifyStream)
	app.POST("/stream-over", streamController.EndStream)
	authGroup := app.Group("auth")
	authGroup.Use(authController.AuthMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/verify-token", authController.VerifyToken)
	}
	streamGroup := app.Group("streams")
	streamGroup.Use(authController.AuthMiddleware.MiddlewareFunc())
	{
		streamGroup.GET("/", streamController.GetStreams)
		streamGroup.POST("/create-key", streamController.CreateKey)
		streamGroup.DELETE("/:id", streamController.DeleteStream)
	}
	nginxGroup := app.Group("/nginx")
	nginxGroup.Use(authController.AuthMiddleware.MiddlewareFunc())
	{
		nginxGroup.GET("/config", confController.GetConfiguration)
		nginxGroup.POST("/config", confController.UpdateConfiguration)
	}
	app.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
