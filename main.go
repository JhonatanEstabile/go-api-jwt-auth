package main

import (
	"api-jwt-auth/controller"
	"api-jwt-auth/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func main() {
	router.POST("/login", controller.Login)
	router.POST("/todo", middleware.TokenAuthMiddleware(), controller.CreateTodo)
	router.POST("/logout", middleware.TokenAuthMiddleware(), controller.Logout)

	log.Fatal(router.Run(":8080"))
}
