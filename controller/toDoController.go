package controller

import (
	"api-jwt-auth/model"
	"api-jwt-auth/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

//CreateTodo register a new to do in db
func CreateTodo(c *gin.Context) {
	var td *model.Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	tokenString := util.ExtractToken(c.Request)

	tokenAuth, err := util.ExtractTokenMetadata(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := util.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	td.UserID = userID

	c.JSON(http.StatusCreated, td)
}
