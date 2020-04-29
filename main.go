package main

import (
	"api-jwt-auth/controller"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func main() {
	router.POST("/login", controller.Login)
	log.Fatal(router.Run(":8080"))
}
